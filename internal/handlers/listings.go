package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/kalpeshWani222/olx-api/internal/httpx"
	"github.com/kalpeshWani222/olx-api/internal/middleware"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandler struct {
	db *sql.DB
	logger *slog.Logger
}

//constructor
// Note => We return the pointer on this contructor because every time when 
//         its using its return the address of it so its does not return the new copy every time 
func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler {
		db : db,
		logger: logger,
	}
}

func (lh ListingHandler) List (w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		rows,err := lh.db.QueryContext(ctx,`
		SELECT id, title, description, price, city, created_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT 100
		`)

		if err != nil {
			lh.logger.Error("listings query error","err",err)
			httpx.Error(w,http.StatusInternalServerError,"Something Went Wrong",httpx.CodeInternalError)
			return
		}

		// Release resources after reading the results
		defer rows.Close()

		listings := []listing{}

		for rows.Next() {
			var l listing
			if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City,&l.CreatedAt); err != nil {
				lh.logger.Error("listings rows scan err","err",err)
				httpx.Error(w,http.StatusInternalServerError,"Something Went Wrong",httpx.CodeInternalError)
				return
			}

		   listings = append(listings, l)
		}

		if err := rows.Err(); err != nil {
			lh.logger.Error("listings rows.err","err",err)
			httpx.Error(w,http.StatusInternalServerError,"Something Went Wrong",httpx.CodeInternalError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		 _ = json.NewEncoder(w).Encode(listings)

    
}
	

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		//receiving the value from the middleware through the context with the use of the exported func
		requestId := middleware.RequestIdFromContext(ctx)
	
		// requestId := ctx.Value(requestId)
		id := r.PathValue("id")
		
		lh.logger.Info("debug log", "listing_id",id)

		_, err := lh.db.ExecContext(ctx,`
		DELETE FROM listings WHERE id = $1`,id)

		if err != nil {
			lh.logger.Error("delete failed", "listing_id",id,"requestId",requestId,"err",err)
			httpx.Error(w,http.StatusInternalServerError,"Something Went Wrong",httpx.CodeInternalError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	
}
