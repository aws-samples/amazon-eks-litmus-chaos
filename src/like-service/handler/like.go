package handler

import (
	"io"
	"like-service/common"
	"like-service/repository"
	"like-service/utils"
	"log"
	"net/http"
	"strconv"
)

type Like struct {
	logger   *log.Logger
	database *repository.Database
	locker   *utils.Database
}

func NewLike(lg *log.Logger, d *repository.Database, lk *utils.Database) *Like {
	return &Like{logger: lg,
		database: d,
		locker:   lk}
}

func GetHealthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Healthy")
}

func (l *Like) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		l.getLikes(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		l.addLike(rw, r)
		return
	}

	// catch all other http verb with 405
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (l *Like) getLikes(rw http.ResponseWriter, r *http.Request) {
	l.logger.Printf("Handle %s %s", r.Method, r.URL)

	likes, err := l.database.FindAllLikes()
	if err != nil {
		http.Error(rw, "Error getting likes", http.StatusInternalServerError)
	}

	likesResp := LikesResponse{
		Likes:    likes,
		Hostname: common.Hostname,
	}

	err = likesResp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (l *Like) addLike(rw http.ResponseWriter, r *http.Request) {
	l.logger.Printf("Handle %s %s", r.Method, r.URL)

	like := LikeRequest{}
	err := like.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	// GET Mutex name from unique id
	mutexName := strconv.Itoa(like.Id)
	mutex := l.locker.Redsync.NewMutex(mutexName)

	// Obtain a lock for our given mutex.
	if err := mutex.Lock(); err != nil {
		panic(err)
	}
	l.logger.Printf("Aquired lock...")

	// Update item in DB
	err = l.database.AddLike(like.Id)
	if err != nil {
		http.Error(rw, "Error encountered executing raw SQL statement", http.StatusInternalServerError)
	}

	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
	l.logger.Printf("Released lock...")

	//	Read likes
	likes, err := l.database.FindAllLikes()
	if err != nil {
		http.Error(rw, "Error getting likes", http.StatusInternalServerError)
	}

	likesResp := LikesResponse{
		Likes:    likes,
		Hostname: common.Hostname,
	}

	err = likesResp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
