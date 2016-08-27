// Copyright © 2016 Dan Esparza <esparza.dan@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/danesparza/tvdb"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	showName               string
	processedDirectoryFlag string
)

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process the files using the passed show name",
	Long: `Processes the files using the passed show name.

Example:
Not available

`,
	Run: func(cmd *cobra.Command, args []string) {
		processwithTVDB()
	},
}

func init() {
	RootCmd.AddCommand(processCmd)

	// Define our flags
	processCmd.Flags().StringVarP(&showName, "show", "s", "", "TV show name to use")
	processCmd.Flags().StringVarP(&processedDirectoryFlag, "directory", "d", "c:\\temp\\", "Directory to move processd files to")
}

func processwithTVDB() {

	//	If we didn't get a show name, exit and indiate we should pass one:
	if showName == "" {
		log.Println("[ERROR] Please pass a show to search with using the '--show' parameter")
		return
	}

	//	Resolve our temporary directory:
	processedDirectory, err := filepath.Abs(filepath.Dir(processedDirectoryFlag))
	if err != nil {
		panic(err)
	}
	log.Printf("[INFO] Processed directory to use: '%v'\n", processedDirectory)

	//	Create our client:
	client := tvdb.TVDBClient{}

	//	Get all series
	request := tvdb.SearchRequest{
		Name: showName}

	seriesMatches, err := client.SeriesSearch(request)
	if err != nil {
		panic(err)
	}
	log.Printf("[INFO] Found %d series matches for '%v':\n", len(seriesMatches), showName)

	for i := 0; i < len(seriesMatches); i++ {
		log.Printf("[INFO] --- Series %d: '%v'\n", i, seriesMatches[i].SeriesName)
	}

	//	Get all episodes for a series:
	episoderequest := tvdb.EpisodeRequest{
		SeriesId: seriesMatches[0].Id}

	episodeMatches, err := client.EpisodesForSeries(episoderequest)
	if err != nil {
		panic(err)
	}
	log.Printf("[INFO] Found %d episodes for '%v' (first match).\n", len(episodeMatches), seriesMatches[0].SeriesName)

	//	Load up the map
	episodes := make(map[string]tvdb.EpisodeResponse)
	for _, episode := range episodeMatches {
		episodes[episode.EpisodeName] = episode
	}

	//	Look for all videos in the current path:
	searchDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	log.Printf("[INFO] Looking for files in '%v'\n", searchDir)

	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {

		//	If it's not a directory...
		if !f.IsDir() {

			filename := f.Name()

			//	If the extension indicates it's a video file:
			videoExtensions := map[string]int{".mp4": 0, ".avi": 0, ".mkv": 0, ".m4v": 0, ".asf": 0, ".mpeg": 0}

			extension := filepath.Ext(filename)
			if _, ok := videoExtensions[extension]; ok {

				//	Get the name to search for:
				name := strings.TrimSuffix(filename, extension)

				//	Try to find it in our episode map:
				if episode, ok := episodes[name]; ok {
					//	We have a match -- describe what we found
					log.Printf("[INFO] Found matching episode for '%v': s%ve%v\n", name, episode.AiredSeason, episode.AiredEpisodeNumber)

					//	Prepare to move the file
					sourcePath := path
					destPath := filepath.Join(processedDirectory, fmt.Sprintf("s%ve%v%v", episode.AiredSeason, episode.AiredEpisodeNumber, extension))

					log.Printf("[INFO] Moving '%v' to '%v'", sourcePath, destPath)
					os.Rename(sourcePath, destPath)
				}
			}
		}

		return nil
	})

}
