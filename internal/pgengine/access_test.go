package pgengine_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cybertec-postgresql/pg_timetable/internal/pgengine"
	"github.com/stretchr/testify/assert"
)

func TestTryLockClientName(t *testing.T) {
	initmockdb(t)
	pgengine.ConfigDb = xdb
	defer db.Close()
	pgengine.ClientName = "pgengine_unit_test"
	pgengine.VerboseLogLevel = false

	t.Run("Check TryLockClientName if everything fine", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectQuery("SELECT pg_try_advisory_lock").WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"}).AddRow(true))
		assert.True(t, pgengine.TryLockClientName(ctx))
	})

	t.Run("Check TryLockClientName if sql fails", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectQuery("SELECT pg_try_advisory_lock").WillReturnError(errors.New("error"))
		assert.False(t, pgengine.TryLockClientName(ctx))
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteChainConfig(t *testing.T) {
	initmockdb(t)
	pgengine.ConfigDb = xdb
	defer db.Close()

	pgengine.VerboseLogLevel = false

	t.Run("Check DeleteChainConfig if everyhing fine", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectExec("INSERT INTO timetable\\.log").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM timetable\\.chain_execution_config").WillReturnResult(sqlmock.NewResult(0, 1))
		assert.True(t, pgengine.DeleteChainConfig(ctx, 0))
	})

	t.Run("Check DeleteChainConfig if sql fails", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectExec("INSERT INTO timetable\\.log").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM timetable\\.chain_execution_config").WillReturnError(errors.New("error"))
		mock.ExpectExec("INSERT INTO timetable\\.log").WillReturnResult(sqlmock.NewResult(0, 1))
		assert.False(t, pgengine.DeleteChainConfig(ctx, 0))
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFixSchedulerCrash(t *testing.T) {
	initmockdb(t)
	pgengine.ConfigDb = xdb
	defer db.Close()

	pgengine.VerboseLogLevel = false

	t.Run("Check FixSchedulerCrash if sql fails", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectExec("INSERT INTO timetable\\.run_status").WillReturnError(errors.New("error"))
		mock.ExpectExec("INSERT INTO timetable\\.log").WillReturnResult(sqlmock.NewResult(0, 1))
		pgengine.FixSchedulerCrash(ctx)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCanProceedChainExecution(t *testing.T) {
	initmockdb(t)
	pgengine.ConfigDb = xdb
	defer db.Close()

	pgengine.VerboseLogLevel = false

	t.Run("Check CanProceedChainExecution if everything fine", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		assert.False(t, pgengine.CanProceedChainExecution(ctx, 0, 0), "Proc count is less than maxinstances")
	})

	t.Run("Check CanProceedChainExecution gets ErrNoRows", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectQuery("SELECT count").WillReturnError(sql.ErrNoRows)
		assert.True(t, pgengine.CanProceedChainExecution(ctx, 0, 0))
	})

	t.Run("Check CanProceedChainExecution if sql fails", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectQuery("SELECT count").WillReturnError(errors.New("error"))
		mock.ExpectExec("INSERT INTO timetable\\.log").WillReturnResult(sqlmock.NewResult(0, 1))
		assert.False(t, pgengine.CanProceedChainExecution(ctx, 0, 0))
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertChainRunStatus(t *testing.T) {
	initmockdb(t)
	pgengine.ConfigDb = xdb
	defer db.Close()

	pgengine.VerboseLogLevel = false

	t.Run("Check InsertChainRunStatus if sql fails", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectQuery("INSERT INTO timetable\\.run_status").WillReturnError(errors.New("error"))
		mock.ExpectExec("INSERT INTO timetable\\.log").WillReturnResult(sqlmock.NewResult(0, 1))
		pgengine.InsertChainRunStatus(ctx, 0, 0)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateChainRunStatus(t *testing.T) {
	initmockdb(t)
	pgengine.ConfigDb = xdb
	defer db.Close()

	pgengine.VerboseLogLevel = false

	t.Run("Check UpdateChainRunStatus if sql fails", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), pgengine.WaitTime*time.Second+2)
		defer cancel()
		mock.ExpectExec("INSERT INTO timetable\\.run_status").WillReturnError(errors.New("error"))
		mock.ExpectExec("INSERT INTO timetable\\.log").WillReturnResult(sqlmock.NewResult(0, 1))
		pgengine.UpdateChainRunStatus(ctx, &pgengine.ChainElementExecution{}, 0, "STATUS")
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestIsAlive(t *testing.T) {
	initmockdb(t)
	assert.False(t, pgengine.IsAlive())

	pgengine.ConfigDb = xdb
	defer db.Close()
	mock.ExpectPing()
	assert.True(t, pgengine.IsAlive())
}
