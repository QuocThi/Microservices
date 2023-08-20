package api

import (
	"context"

	"github.com/go-logr/logr"

	db "log-service/internal/database"
	logs "log-service/proto"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	DB  db.DB
	Log logr.Logger
}

func NewGRPCLogServer(db db.DB, l logr.Logger) *LogServer {
	return &LogServer{
		UnimplementedLogServiceServer: logs.UnimplementedLogServiceServer{},
		DB:                            db,
		Log:                           l,
	}
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := db.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.DB.Insert(logEntry)
	if err != nil {
		l.Log.Error(err, "failed to insert db")
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// return response
	res := &logs.LogResponse{Result: "logged!"}
	return res, nil
}
