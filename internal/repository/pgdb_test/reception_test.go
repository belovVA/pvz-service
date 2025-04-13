package pgdb_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"pvz-service/internal/repository/pgdb"
	modelRepo "pvz-service/internal/repository/pgdb/model"
)

func TestCreateReception(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := pgdb.NewReceptionRepository(mock)

	t.Run("success", func(t *testing.T) {
		pvzID := uuid.New()
		expectedID := uuid.New()

		mock.ExpectQuery("INSERT INTO reception").
			WithArgs(pvzID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(expectedID))

		id, err := repo.CreateReception(context.Background(), pvzID)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		pvzID := uuid.New()

		mock.ExpectQuery("INSERT INTO reception").
			WithArgs(pvzID).
			WillReturnError(errors.New("database error"))

		id, err := repo.CreateReception(context.Background(), pvzID)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetReceptionByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := pgdb.NewReceptionRepository(mock)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		expectedReception := modelRepo.Reception{
			ID:             id,
			DateTime:       time.Now(),
			IsClosedStatus: false,
			PvzID:          uuid.New(),
		}

		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WithArgs(id.String()).
			WillReturnRows(pgxmock.NewRows([]string{"id", "date_time", "is_closed", "pvz_id"}).
				AddRow(expectedReception.ID, expectedReception.DateTime, expectedReception.IsClosedStatus, expectedReception.PvzID))

		reception, err := repo.GetReceptionByID(context.Background(), id)

		assert.NoError(t, err)
		assert.Equal(t, id, reception.ID)
		assert.Equal(t, expectedReception.DateTime, reception.DateTime)
		assert.Equal(t, expectedReception.IsClosedStatus, reception.IsClosed)
		assert.Equal(t, expectedReception.PvzID, reception.PvzID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		id := uuid.New()

		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WithArgs(id.String()).
			WillReturnError(errors.New("no rows in result set"))

		reception, err := repo.GetReceptionByID(context.Background(), id)

		assert.Error(t, err)
		assert.Nil(t, reception)
		assert.Equal(t, "reception not found", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetLastReception(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := pgdb.NewReceptionRepository(mock)

	t.Run("success", func(t *testing.T) {
		pvzID := uuid.New()
		expectedReception := modelRepo.Reception{
			ID:             uuid.New(),
			DateTime:       time.Now(),
			IsClosedStatus: false,
			PvzID:          pvzID,
		}

		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WithArgs(pvzID.String()).
			WillReturnRows(pgxmock.NewRows([]string{"id", "date_time", "is_closed", "pvz_id"}).
				AddRow(expectedReception.ID, expectedReception.DateTime, expectedReception.IsClosedStatus, expectedReception.PvzID))

		reception, err := repo.GetLastReception(context.Background(), pvzID)

		assert.NoError(t, err)
		assert.Equal(t, expectedReception.ID, reception.ID)
		assert.Equal(t, expectedReception.DateTime, reception.DateTime)
		assert.Equal(t, expectedReception.IsClosedStatus, reception.IsClosed)
		assert.Equal(t, expectedReception.PvzID, reception.PvzID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		pvzID := uuid.New()

		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WithArgs(pvzID.String()).
			WillReturnError(errors.New("no rows in result set"))

		reception, err := repo.GetLastReception(context.Background(), pvzID)

		assert.Error(t, err)
		assert.Nil(t, reception)
		assert.Equal(t, "reception not found", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCloseReception(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := pgdb.NewReceptionRepository(mock)

	t.Run("success", func(t *testing.T) {
		receptionID := uuid.New()

		mock.ExpectExec("UPDATE reception").
			WithArgs(true, receptionID.String()).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.CloseReception(context.Background(), receptionID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows affected", func(t *testing.T) {
		receptionID := uuid.New()

		mock.ExpectExec("UPDATE reception").
			WithArgs(true, receptionID.String()).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := repo.CloseReception(context.Background(), receptionID)

		assert.Error(t, err)
		assert.Equal(t, "no rows affected", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		receptionID := uuid.New()

		mock.ExpectExec("UPDATE reception").
			WithArgs(true, receptionID.String()).
			WillReturnError(errors.New("database error"))

		err := repo.CloseReception(context.Background(), receptionID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute query")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetReceptionsSliceWithTimeRange(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := pgdb.NewReceptionRepository(mock)

	now := time.Now()
	begin := now.Add(-24 * time.Hour)
	end := now.Add(24 * time.Hour)

	t.Run("success with time range", func(t *testing.T) {
		expectedReceptions := []modelRepo.Reception{
			{
				ID:             uuid.New(),
				DateTime:       now.Add(-12 * time.Hour),
				IsClosedStatus: false,
				PvzID:          uuid.New(),
			},
			{
				ID:             uuid.New(),
				DateTime:       now.Add(12 * time.Hour),
				IsClosedStatus: true,
				PvzID:          uuid.New(),
			},
		}

		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WithArgs(begin, end).
			WillReturnRows(pgxmock.NewRows([]string{"id", "date_time", "is_closed", "pvz_id"}).
				AddRow(expectedReceptions[0].ID, expectedReceptions[0].DateTime, expectedReceptions[0].IsClosedStatus, expectedReceptions[0].PvzID).
				AddRow(expectedReceptions[1].ID, expectedReceptions[1].DateTime, expectedReceptions[1].IsClosedStatus, expectedReceptions[1].PvzID))

		receptions, err := repo.GetReceptionsSliceWithTimeRange(context.Background(), begin, end)

		assert.NoError(t, err)
		assert.Len(t, receptions, 2)
		assert.Equal(t, expectedReceptions[0].ID, receptions[0].ID)
		assert.Equal(t, expectedReceptions[1].ID, receptions[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success with begin only", func(t *testing.T) {
		expectedReception := modelRepo.Reception{
			ID:             uuid.New(),
			DateTime:       now.Add(12 * time.Hour),
			IsClosedStatus: true,
			PvzID:          uuid.New(),
		}

		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WithArgs(begin).
			WillReturnRows(pgxmock.NewRows([]string{"id", "date_time", "is_closed", "pvz_id"}).
				AddRow(expectedReception.ID, expectedReception.DateTime, expectedReception.IsClosedStatus, expectedReception.PvzID))

		receptions, err := repo.GetReceptionsSliceWithTimeRange(context.Background(), begin, time.Time{})

		assert.NoError(t, err)
		assert.Len(t, receptions, 1)
		assert.Equal(t, expectedReception.ID, receptions[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success with end only", func(t *testing.T) {
		expectedReception := modelRepo.Reception{
			ID:             uuid.New(),
			DateTime:       now.Add(-12 * time.Hour),
			IsClosedStatus: false,
			PvzID:          uuid.New(),
		}

		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WithArgs(end).
			WillReturnRows(pgxmock.NewRows([]string{"id", "date_time", "is_closed", "pvz_id"}).
				AddRow(expectedReception.ID, expectedReception.DateTime, expectedReception.IsClosedStatus, expectedReception.PvzID))

		receptions, err := repo.GetReceptionsSliceWithTimeRange(context.Background(), time.Time{}, end)

		assert.NoError(t, err)
		assert.Len(t, receptions, 1)
		assert.Equal(t, expectedReception.ID, receptions[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without time range", func(t *testing.T) {
		expectedReceptions := []modelRepo.Reception{
			{
				ID:             uuid.New(),
				DateTime:       now.Add(-48 * time.Hour),
				IsClosedStatus: false,
				PvzID:          uuid.New(),
			},
			{
				ID:             uuid.New(),
				DateTime:       now.Add(48 * time.Hour),
				IsClosedStatus: true,
				PvzID:          uuid.New(),
			},
		}

		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WillReturnRows(pgxmock.NewRows([]string{"id", "date_time", "is_closed", "pvz_id"}).
				AddRow(expectedReceptions[0].ID, expectedReceptions[0].DateTime, expectedReceptions[0].IsClosedStatus, expectedReceptions[0].PvzID).
				AddRow(expectedReceptions[1].ID, expectedReceptions[1].DateTime, expectedReceptions[1].IsClosedStatus, expectedReceptions[1].PvzID))

		receptions, err := repo.GetReceptionsSliceWithTimeRange(context.Background(), time.Time{}, time.Time{})

		assert.NoError(t, err)
		assert.Len(t, receptions, 2)
		assert.Equal(t, expectedReceptions[0].ID, receptions[0].ID)
		assert.Equal(t, expectedReceptions[1].ID, receptions[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WithArgs(begin, end).
			WillReturnError(errors.New("database error"))

		receptions, err := repo.GetReceptionsSliceWithTimeRange(context.Background(), begin, end)

		assert.Error(t, err)
		assert.Nil(t, receptions)
		assert.Contains(t, err.Error(), "failed to execute query")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, date_time, is_closed, pvz_id FROM reception").
			WithArgs(begin, end).
			WillReturnRows(pgxmock.NewRows([]string{"id", "date_time", "is_closed", "pvz_id"}).
				AddRow("invalid-uuid", now, false, uuid.New())) // Invalid UUID format

		receptions, err := repo.GetReceptionsSliceWithTimeRange(context.Background(), begin, end)

		assert.Error(t, err)
		assert.Nil(t, receptions)
		assert.Contains(t, err.Error(), "failed to scan row")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
