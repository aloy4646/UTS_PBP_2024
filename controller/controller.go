package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"uts/model"

	"github.com/gorilla/mux"
)

func GetAllRoomsByGame(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id_game := vars["id_game"]

	query := "SELECT id, room_name from rooms WHERE id_game=" + id_game

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		responseMessage(w, 500, "Internal server error")
		return
	}

	var room model.RoomWithoutGame
	var rooms []model.RoomWithoutGame
	for rows.Next() {
		if err := rows.Scan(&room.Id, &room.Room_name); err != nil {
			log.Println(err.Error())
			responseMessage(w, 170, "Data error")
			return
		} else {
			rooms = append(rooms, room)
			if err != nil {
				log.Println(err)
				responseMessage(w, 500, "Internal server error")
				return
			}
		}
	}

	sendSuccessResponseWithData(w, rooms)
}

func GetDetailRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	query := "SELECT r.id, r.room_name, p.id, a.id, a.username from rooms r JOIN participants p ON r.id = p.id_room JOIN accounts a ON p.id_account = a.id WHERE r.id=" + id

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		responseMessage(w, 500, "Internal server error")
		return
	}

	var room model.Room
	for rows.Next() {
		var participantDetail model.ParticipantDetail
		if err := rows.Scan(&room.Id, &room.Room_name, &participantDetail.Id, &participantDetail.Id_account, &participantDetail.Username); err != nil {
			log.Println(err.Error())
			responseMessage(w, 170, "Data error")
			return
		} else {
			room.Participants = append(room.Participants, participantDetail)
			if err != nil {
				log.Println(err)
				responseMessage(w, 500, "Internal server error")
				return
			}
		}
	}

	//tambahan supaya jsonnya sesuai dengan yang disoal
	var responseRoomDetail model.ResponseRoomDetail
	responseRoomDetail.Room = room

	sendSuccessResponseWithData(w, responseRoomDetail)
}

// Insert room --> insert accoount sebagai sebuah participant disebuah game
// bukan nambah data di tabel room tetapi nambah data di tabel participants
func InsertRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		responseMessage(w, 500, "Internal server error")
		log.Println(err)
		return
	}

	id_room, _ := strconv.Atoi(r.Form.Get("id_room"))
	id_account, _ := strconv.Atoi(r.Form.Get("id_account"))

	// cek apakah sebuah account sudah terdaftar di room tersebut atau belum
	// *opsional
	if cekRedundansi(id_room, id_account, db) == 1 {
		responseMessage(w, 422, "Insert Gagal, player sudah ada")
		return
	}

	// cek player pada room apakah sudah max atau belum
	if checkRoomPlayer(id_room, db) {
		responseMessage(w, 422, "Insert Gagal, jumlah player telah mencapai batas")
		return
	}

	_, errQuery := db.Exec("INSERT INTO participants(id_room, id_account) VALUES (?,?)",
		id_room,
		id_account,
	)

	if errQuery != nil {
		responseMessage(w, 500, "Internal server error")
		log.Println(errQuery)

		return
	}

	responseMessage(w, 200, "Insert Success")
}

func cekRedundansi(id_room int, id_account int, db *sql.DB) int {
	var jumlahData int
	err := db.QueryRow("SELECT COUNT(*) FROM participants WHERE id_room =? AND id_account=?", id_room, id_account).Scan(&jumlahData)
	if err != nil {
		log.Println("Error saat menghitung jumlah pemain:", err)
	}

	if jumlahData == 1 {
		return 1
	} else {
		return 0
	}
}

func checkRoomPlayer(id_room int, db *sql.DB) bool {
	max := true

	// cek jumlah player saat ini
	var jumlahPlayerSaatIni int
	err := db.QueryRow("SELECT COUNT(*) FROM participants WHERE id_room = ?", id_room).Scan(&jumlahPlayerSaatIni)
	if err != nil {
		log.Println("Error saat menghitung jumlah pemain:", err)
	}

	// get max player dari game
	rows, err := db.Query("SELECT g.max_player FROM games g JOIN rooms r ON g.id = r.id_game WHERE r.id=? LIMIT 1", id_room)
	if err != nil {
		log.Println("Error saat get max player:", err)
		return max
	}
	var max_player int

	for rows.Next() {
		if errScan := rows.Scan(&max_player); errScan != nil {
			if errScan != nil {
				log.Println("Error saat scan:", errScan)
				return max
			}
		}
	}

	if jumlahPlayerSaatIni < max_player {
		max = false
	}

	return max
}

func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		responseMessage(w, 500, "Internal server error")
		log.Println(err)
		return
	}

	id_room := r.URL.Query().Get("id_room")
	id_account := r.URL.Query().Get("id_account")

	resultQuery, errQuery := db.Exec("DELETE FROM participants WHERE id_room=? AND id_account=?", id_room, id_account)
	if errQuery != nil {
		responseMessage(w, 500, "Internal server error")
		log.Println(errQuery)
		return
	}

	// query DELETE akan tetap berhasil meskipun 0 data terhapus
	// cek apakah ada data yang terhapus
	rowsAffected, _ := resultQuery.RowsAffected()

	if rowsAffected > 0 {
		responseMessage(w, 200, "Delete Success")
	} else {
		responseMessage(w, 407, "Failed, 0 rows affected")
	}
}
