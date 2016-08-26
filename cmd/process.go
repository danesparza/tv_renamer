// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"github.com/danesparza/tvdb"
	"github.com/spf13/cobra"
	"log"
)

var (
	showName string
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
		// TODO: Work your own magic here
		testTVDB()
	},
}

func init() {
	RootCmd.AddCommand(processCmd)

	// Define our flags
	processCmd.Flags().StringVarP(&showName, "show", "s", "", "TV show name to use")

}

func testTVDB() {

	//	If we didn't get a show name, exit and indiate we should pass one:
	if showName == "" {
		log.Println("[ERROR] Please pass a show to search with using the '--show' parameter")
		return
	}

	//	Create our client:
	client := tvdb.TVDBClient{}

	//	Get all series
	request := tvdb.SearchRequest{
		Name: showName}

	seriesMatches, err := client.SeriesSearch(request)
	if err != nil {
		panic(err)
	}
	log.Printf("[INFO] Found '%d' series matches for '%v'\n", len(seriesMatches), showName)

	//	Get all episodes for a series:
	episoderequest := tvdb.EpisodeRequest{
		SeriesId: seriesMatches[0].Id}

	episodeMatches, err := client.EpisodesForSeries(episoderequest)
	if err != nil {
		panic(err)
	}
	log.Printf("[INFO] Found %d episodes for '%v'.\n", len(episodeMatches), seriesMatches[0].SeriesName)

	log.Printf("[INFO] First season: %v First episode: %v \n", episodeMatches[0].AiredSeason, episodeMatches[0].AiredEpisodeNumber)

}
