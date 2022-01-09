package access_logger

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/talon-one/assignment-props/dat"
	"go.uber.org/zap"
	"sync"
)

var AccessLogChan chan LogData

var RequestResponseMap sync.Map

type LogData interface {
	GetUUID() uuid.UUID
	GetBody() string
}

type RequestLogData struct {
	UUID uuid.UUID
	Body string
}

func (r RequestLogData) GetUUID() uuid.UUID {
	return r.UUID
}

func (r RequestLogData) GetBody() string {
	return r.Body
}

type ResponseLogData struct {
	UUID uuid.UUID
	Body string
}

func (r ResponseLogData) GetUUID() uuid.UUID {
	return r.UUID
}

func (r ResponseLogData) GetBody() string {
	return r.Body
}

func Init() {
	AccessLogChan = make(chan LogData, 100)
}

func Run(logger *zap.Logger) {
	select {
	case msg := <-AccessLogChan:
		msgUuid := msg.GetUUID()
		dataContext := dataContext(msgUuid, logger)
		switch msg.(type) {
		case RequestLogData:
			err := dataContext.AccessLog.InsertRequest(msg.GetBody())
			if err != nil {
				rollback(msgUuid, dataContext)
				logger.Fatal("InsertRequest failed", zap.Error(err))
			}

		case ResponseLogData:
			err := dataContext.AccessLog.InsertResponse(msg.GetBody())
			if err != nil {
				rollback(msgUuid, dataContext)
				logger.Fatal("InsertRequest failed", zap.Error(err))
			}
		}
		commitOrSave(msgUuid, dataContext)
	}
}

func dataContext(uuid uuid.UUID, logger *zap.Logger) *dat.Context {
	dc, ok := RequestResponseMap.Load(uuid)
	if ok {
		dataContext := dc.(dat.Context)
		return &dataContext
	}
	dataContext, err := dat.NewDataContext(context.Background())
	if err != nil {
		logger.Fatal("Data context init failed", zap.Error(err))
	}
	return dataContext
}

func commitOrSave(uuid uuid.UUID, dataContext *dat.Context) {
	dc, ok := RequestResponseMap.Load(uuid)
	if ok {
		dataContext := dc.(dat.Context)
		dataContext.Commit()
		RequestResponseMap.Delete(uuid)
	} else {
		RequestResponseMap.Store(uuid, dataContext)
	}
}

func rollback(uuid uuid.UUID, dataContext *dat.Context) {
	dc, ok := RequestResponseMap.Load(uuid)
	if ok {
		dataContext := dc.(dat.Context)
		dataContext.Rollback()
		RequestResponseMap.Delete(uuid)
	} else {
		dataContext.Rollback()
	}
}