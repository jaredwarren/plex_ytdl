package server

import (
	"testing"
)

func TestListSongHandler(t *testing.T) {
	// db := initDB(t)

	// err := db.UpdateSong(&model.Song{
	// 	ID:   "test_song",
	// 	RFID: "test_song_rfid",
	// })
	// assert.NoError(t, err)

	// // sanity test
	// songs, err := db.ListSongs()
	// assert.Nil(t, err)
	// assert.Len(t, songs, 1)
	// assert.Equal(t, "test_song", songs[0].ID)

	// s := &Server{
	// 	db:     db,
	// 	logger: log.NewNoOpLogger(),
	// 	downloader: &downloader.MockDownloader{
	// 		Response: map[string]*youtube.Video{"new url": {
	// 			ID:    "test_song",
	// 			Title: "song title",
	// 			Thumbnails: youtube.Thumbnails{{
	// 				URL: "thumb_url",
	// 			}},
	// 		}},
	// 	},
	// }

	// { // test List song
	// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// 	w := httptest.NewRecorder()

	// 	s.ListSongHandler(w, req)

	// 	res := w.Result()
	// 	defer res.Body.Close()
	// 	data, err := ioutil.ReadAll(res.Body)
	// 	assert.NoError(t, err)

	// 	assert.Contains(t, string(data), `test_song_rfid`)
	// }

	// { // test Edit song form
	// 	req := httptest.NewRequest(http.MethodGet, "/song/test_song", nil)
	// 	req = mux.SetURLVars(req, map[string]string{
	// 		"song_id": "test_song",
	// 	})
	// 	w := httptest.NewRecorder()

	// 	s.EditSongFormHandler(w, req)

	// 	res := w.Result()
	// 	defer res.Body.Close()
	// 	data, err := ioutil.ReadAll(res.Body)
	// 	assert.NoError(t, err)

	// 	assert.Contains(t, string(data), `<form action="/song/test_song" method="post">`)
	// }

	// { // test Post Edit
	// 	form := url.Values{}
	// 	form.Add("url", "new url")
	// 	form.Add("rfid", "new:rfid")
	// 	req := httptest.NewRequest(http.MethodPost, "/song/test_song", strings.NewReader(form.Encode()))
	// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// 	req = mux.SetURLVars(req, map[string]string{
	// 		"song_id": "test_song",
	// 	})
	// 	w := httptest.NewRecorder()

	// 	s.UpdateSongHandler(w, req)

	// 	res := w.Result()
	// 	assert.Equal(t, http.StatusMovedPermanently, res.StatusCode)
	// 	defer res.Body.Close()
	// 	assert.NoError(t, err)

	// 	{ // that that list songs has new rfid
	// 		req := httptest.NewRequest(http.MethodGet, "/", nil)
	// 		w := httptest.NewRecorder()

	// 		s.ListSongHandler(w, req)

	// 		res := w.Result()
	// 		defer res.Body.Close()
	// 		data, err := ioutil.ReadAll(res.Body)
	// 		assert.NoError(t, err)

	// 		assert.Contains(t, string(data), "newrfid")
	// 	}
	// }
}
