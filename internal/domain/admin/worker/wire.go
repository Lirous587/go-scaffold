package worker

import (
	"go.uber.org/zap"
	"scaffold/internal/domain/admin/infrastructure/cache"
	"scaffold/internal/domain/admin/infrastructure/db"
	"time"
)

type Worker interface {
	Start()
	Stop()
}

type worker struct {
	db       db.DB
	cache    cache.Cache
	ticker   *time.Ticker
	stopChan chan struct{}
}

func NewWorker(db db.DB, cache cache.Cache) Worker {
	return &worker{
		db:       db,
		cache:    cache,
		stopChan: make(chan struct{}),
	}
}

const workerInternal = 24 * time.Hour

func (w *worker) Start() {
	w.ticker = time.NewTicker(workerInternal)

	go func() {
		// 启动后立即执行一次同步
		if err := w.syncTasks(); err != nil {
			sugar := zap.L().Sugar()
			sugar.Errorf("同步任务失败: %+v", err)
		}

		for {
			select {
			case <-w.ticker.C:
				if err := w.syncTasks(); err != nil {
					zap.L().Error("定时同步任务失败: %v", zap.Error(err))
				}
			case <-w.stopChan:
				w.ticker.Stop()
				return
			}
		}
	}()
	zap.L().Info("访问次数同步任务已启动")
}

func (w *worker) Stop() {
	close(w.stopChan)
	zap.L().Info("访问次数同步任务已停止")
}

func (w *worker) syncTasks() error {
	mockFunc := func() {
		var a = 1
		a++
	}
	mockFunc()
	// 将普通logger转为Sugar logger
	sugar := zap.L().Sugar()
	sugar.Infof("完成mock任务\n")
	return nil
}
