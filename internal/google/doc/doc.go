package doc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	srv   *docs.Service
	DocId string
)

func init() {
	c, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to read credentials: %v", err))
	}

	config, err := google.ConfigFromJSON(c, "https://www.googleapis.com/auth/documents")
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to get config JSON from googleapis: %v", err))
	}
	client, err := getClient(config)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to get client: %v", err))
	}

	srv, err = docs.New(client)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to new server: %v", err))
	}
}

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err2 := getTokenFromWeb(config)
		if err2 != nil {
			return nil, err2
		}
		if err2 := saveToken(tokFile, tok); err2 != nil {
			return nil, err2
		}
	}
	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return errors.Wrapf(err, "Unable to cache OAuth token")
	}
	json.NewEncoder(f).Encode(token)
	return nil
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, errors.Wrapf(err, "Unable to read authorization code")
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve token from web")
	}
	return tok, nil
}

func InsertText(message string) error {
	doc, err := srv.Documents.Get(DocId).Do()
	if err != nil {
		return errors.Wrapf(err, "Unable to retrieve data from document")
	}

	b := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				InsertText: &docs.InsertTextRequest{
					Text: message,
					Location: &docs.Location{
						Index: doc.Body.Content[len(doc.Body.Content)-1].EndIndex - 1,
					},
				},
			},
		},
	}

	if _, err := srv.Documents.BatchUpdate(DocId, b).Do(); err != nil {
		return err
	}

	return nil
}

func InsertImage(url string) error {
	doc, err := srv.Documents.Get(DocId).Do()
	if err != nil {
		return errors.Wrapf(err, "Unable to retrieve data from document")
	}

	b := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				InsertInlineImage: &docs.InsertInlineImageRequest{
					Location: &docs.Location{
						Index: doc.Body.Content[len(doc.Body.Content)-1].EndIndex - 1,
					},
					ObjectSize: nil,
					Uri:        url,
				},
			},
		},
	}

	if _, err := srv.Documents.BatchUpdate(DocId, b).Do(); err != nil {
		return err
	}

	return nil
}

func InsertHyperLink(message, url string) error {
	doc, err := srv.Documents.Get(DocId).Do()
	if err != nil {
		return errors.Wrapf(err, "Unable to retrieve data from document")
	}

	b := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				InsertText: &docs.InsertTextRequest{
					Text: message + "\n\n",
					Location: &docs.Location{
						Index: doc.Body.Content[len(doc.Body.Content)-1].EndIndex - 1,
					},
				},
			},
		},
	}

	if _, err := srv.Documents.BatchUpdate(DocId, b).Do(); err != nil {
		return err
	}

	b = &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				UpdateTextStyle: &docs.UpdateTextStyleRequest{
					Fields: "*",
					Range: &docs.Range{
						StartIndex: doc.Body.Content[len(doc.Body.Content)-1].EndIndex - 1,
						EndIndex:   doc.Body.Content[len(doc.Body.Content)-1].EndIndex + int64(len(message)),
					},
					TextStyle: &docs.TextStyle{
						Link: &docs.Link{
							Url: url,
						},
						ForegroundColor: &docs.OptionalColor{
							Color: &docs.Color{
								RgbColor: &docs.RgbColor{
									Blue:  0.8,
									Green: 0.33333334,
									Red:   0.06666667,
								},
							},
						},
					},
				},
			},
		},
	}

	if _, err := srv.Documents.BatchUpdate(DocId, b).Do(); err != nil {
		return err
	}

	return nil
}
