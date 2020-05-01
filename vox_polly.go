// wevox_polly.go
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

func textToSpeechPolly(audioID int, audioName string, audioText string, voice string) (WSAudio, WSErrors) {
	var ttsAudio WSAudio
	var ttsError WSErrors

	ttsInput := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		SampleRate: aws.String("22050"),
		Text: aws.String(audioText),
		TextType: aws.String("text"),
		VoiceId: aws.String(voice),
	}

	svc := polly.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))
	fmt.Println("SVC")

	ttsResult, svcErr := svc.SynthesizeSpeech(ttsInput)

	if svcErr != nil {
		if pollyErr, ok := svcErr.(awserr.Error); ok {
			switch pollyErr.Code() {

			case polly.ErrCodeTextLengthExceededException:
				fmt.Println(polly.ErrCodeTextLengthExceededException, pollyErr.Error())

			case polly.ErrCodeInvalidSampleRateException:
				fmt.Println(polly.ErrCodeInvalidSampleRateException, pollyErr.Error())

			case polly.ErrCodeInvalidSsmlException:
				fmt.Println(polly.ErrCodeInvalidSsmlException, pollyErr.Error())

			case polly.ErrCodeLexiconNotFoundException:
				fmt.Println(polly.ErrCodeLexiconNotFoundException, pollyErr.Error())

			case polly.ErrCodeServiceFailureException:
				fmt.Println(polly.ErrCodeServiceFailureException, pollyErr.Error())

			case polly.ErrCodeMarksNotSupportedForFormatException:
				fmt.Println(polly.ErrCodeMarksNotSupportedForFormatException, pollyErr.Error())

			case polly.ErrCodeSsmlMarksNotSupportedForTextTypeException:
				fmt.Println(polly.ErrCodeSsmlMarksNotSupportedForTextTypeException, pollyErr.Error())

			case polly.ErrCodeLanguageNotSupportedException:
				fmt.Println(polly.ErrCodeLanguageNotSupportedException, pollyErr.Error())

			default:
				fmt.Println(pollyErr.Error())
				ttsError.AudioID = audioID
				ttsError.AudioName = audioName
				ttsError.ErrorDescription = pollyErr.Error()
			}
		} else {

			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.

			fmt.Println(svcErr.Error())
			ttsError.AudioID = audioID
			ttsError.AudioName = audioName
			ttsError.ErrorDescription = svcErr.Error()
		}
		return ttsAudio, ttsError
	}

	outFile, err := os.Create(outputPath + audioName + ".mp3")
	defer outFile.Close()

	if err != nil {
		ttsError.AudioID = audioID
		ttsError.AudioName = audioName
		ttsError.ErrorDescription = err.Error()
	}

	outFileSize, err := io.Copy(outFile, ttsResult.AudioStream)
	fmt.Println(ttsResult.String())

	if err != nil {
		ttsError.AudioID = audioID
		ttsError.AudioName = audioName
		ttsError.ErrorDescription = err.Error()
	}

	audioLength := fmt.Sprintf("%.2f", (.000167 * float64(outFileSize)))

	ttsAudio.AudioLengthSec, _ = strconv.ParseFloat(audioLength, 64)
	ttsAudio.AudioID = audioID
	ttsAudio.AudioName = audioName
	ttsAudio.AudioURL = "http://" + downloadServerURL + "/mp3/" + audioName + ".mp3"
	ttsAudio.AudioText = audioText
	ttsAudio.AudioSQL = ""

	fmt.Printf("Audio is %v [%d bytes], length is %.2f seconds\n", ttsAudio.AudioName, outFileSize, ttsAudio.AudioLengthSec)

	return ttsAudio, ttsError
}
