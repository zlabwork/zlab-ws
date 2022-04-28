package seq

import (
	"app/service/repository/mysql"
	"app/service/repository/redis"
	"context"
)

// TODO:: SEQ集群 分布式锁
// @docs https://mp.weixin.qq.com/s/JqIJupVKUNuQYIDDxRtfqA

// 步长机制
// 通过 step 计算 max_seq, 每隔 step 写入一次硬盘, 减少硬盘IO写入
const step = 10000

// 起始ID, 必须大于1
const seqStartId = 100

// 分号段共享存储
// 减少 max_seq 数量, 防止重启时要读取大量的 max_seq 数据加载到内存中
// 1万个人共享一个 max_seq
const sectionSize = 10000

type Service struct {
	ctx context.Context
	rd  *redis.SeqRepository
	db  *mysql.SeqRepository
}

func NewSeqService() (*Service, error) {

	ctx := context.TODO()
	rd, err := redis.NewSeqRepository()
	if err != nil {
		return nil, err
	}
	db, err := mysql.NewSeqRepository()
	if err != nil {
		return nil, err
	}
	return &Service{ctx: ctx, rd: rd, db: db}, nil
}

func (srv *Service) GetSeqId(userId int64) (uint64, error) {

	// 1. get seq_id
	seq, err := srv.rd.GetSeqId(srv.ctx, userId)
	if err != nil {
		return 0, err
	}

	// 2. if restart or first time
	if seq < seqStartId {
		srv.rd.SetSeqId(srv.ctx, userId, seqStartId)
		sec := userId % sectionSize
		max, _ := srv.rd.GetMaxSeq(srv.ctx, sec)
		if max < step {
			srv.rd.SetMaxSeq(srv.ctx, sec, (seqStartId/step+1)*step)
		}
		return seqStartId, err
	}

	// 2.
	if seq%step != 0 {
		return seq, nil
	}

	// 3. get max_seq
	sec := userId % sectionSize
	max, _ := srv.rd.GetMaxSeq(srv.ctx, sec)

	// 4. write to database
	if seq >= max {
		srv.rd.SetMaxSeq(srv.ctx, sec, max+step)
		go func() {
			srv.db.SetMax(sec, max+step)
		}()
	}

	return seq, nil
}

func (srv *Service) LoadSeq() error {

	// 1. get from database
	data, err := srv.db.GetAll()
	if err != nil {
		return err
	}

	// 2. init
	if len(data) == 0 {
		err = srv.db.FillInitData(sectionSize, (seqStartId/step+1)*step)
		if err != nil {
			return err
		}

		// 3. read again
		data, err = srv.db.GetAll()
		if err != nil {
			return err
		}
	}

	// 4. set to cache
	return srv.rd.LoadMaxSeq(srv.ctx, data)
}
