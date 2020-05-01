// wevox_http.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type WSText struct {
	AudioID int `json:"audio_id"`
	AudioName string `json:"audio_name"`
	AudioText string `json:"audio_text"`
}

type WSAudio struct {
	AudioID int `json:"audio_id"`
	AudioName string `json:"audio_name"`
	AudioLengthSec float64 `json:"audio_length_sec"`
	AudioText string `json:"audio_text"`
	AudioURL string `json:"audio_url"`
	AudioSQL string `json:"audio_sql"`
}

type WSErrors struct {
	ErrorDescription string `json:"error_description"`
	AudioID int `json:"audio_id"`
	AudioName string `json:"audio_name"`
}

type TextJSON struct {
	Mode string `json:"mode"`
	Voice string `json:"voice"`
	Text []WSText `json:"text_to_convert"`
}

type AudioJSON struct {
	NumberOfAudios int `json:"number_of_audios"`
	Audio []WSAudio `json:"audio"`
	AudioErrors []WSErrors `json:"audio_errors"`
}

func (audio *AudioJSON) AddAudioEntry(item WSAudio) {
	audio.Audio = append(audio.Audio, item)
}

func (audio *AudioJSON) AddErrorEntry(item WSErrors) {
	audio.AudioErrors = append(audio.AudioErrors, item)
}

func showHTMLHelp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Called showHtmlHelp")
}

func convert(w http.ResponseWriter, r *http.Request) {

	var textJSON TextJSON
	err := r.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &textJSON); err != nil {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}

	}

	fmt.Printf("%d records to batch\n", len(textJSON.Text))
	textJSON.Voice = "Joanna" // remove this line to customize voice through JSON request.
	var audioJSON AudioJSON

	for ctr := 0; ctr < len(textJSON.Text); ctr++ {
		audio, errors := textToSpeechPolly(textJSON.Text[ctr].AudioID, textJSON.Text[ctr].AudioName, textJSON.Text[ctr].AudioText, textJSON.Voice)

		if len(errors.ErrorDescription) != 0 {
			audioJSON.AddErrorEntry(errors)
		} else {
			audioJSON.AddAudioEntry(audio)
		}
		fmt.Println("-----")
	}

	audioJSON.NumberOfAudios = len(audioJSON.Audio)
	returnJSON, err := json.Marshal(audioJSON)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(returnJSON)
}
