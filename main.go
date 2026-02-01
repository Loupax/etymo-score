package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)
type GetCountryArgs struct {
    Country string `json:"country" jsonschema:"The name of the country to fetch variations for"`
}
type GetCountryResult struct {
	Names []string `json:"names"`
}
const Model = "gemini-3-pro-preview"//"gemini-2.5-flash"
func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, Model, &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}
	langTool, err := functiontool.New(functiontool.Config{
		Name: "GetCountryVariations",
		Description:"Fetches the variations a country name has among other countries",
	},GetCountryVariations)
	if err!=nil{
		log.Fatalf("Failed registering language fetching tool: %v", err)
	}
	a, err := llmagent.New(llmagent.Config{
		Name:        "country_name_scorer",
		Model:       model,
		Instruction: "Your goal is to calculate an Etymological Variation Score for a country. " +
		"1. Use the provided tool to fetch all known variations of the country's name. " +
		"2. Group these names by their 'Cognate' (common etymological root). " +
		"   - Names derived from the same linguistic ancestor must be in the same group (e.g., 'Spain', 'Espagne', and 'Spanien' all come from 'Hispania'). " +
		"   - Names from distinct roots must be in separate groups (e.g., 'Hungary' from 'Onogur' vs 'Magyarország' from 'Magyar'). " +
		"3. Count the number of unique etymological groups. " +
		"4. Return the final count as the score. Briefly explain your grouping logic in your final response.",
		Description: "A tool that calculates the language score of counteries",
		Tools: []tool.Tool{
			langTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(a),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}

// Define a tool the agent can call
func GetCountryVariations(ctx tool.Context, args GetCountryArgs) (GetCountryResult, error) {
	// 1. Define the SPARQL query
	// This query finds an item with the label (args.Country), ensures it's an instance of a country (Q6256),
	// and then pulls all labels associated with that Wikidata ID.
	sparqlQuery := fmt.Sprintf(`
		SELECT DISTINCT ?label WHERE {
		  ?country rdfs:label "%s"@en ;
		           wdt:P31 wd:Q6256 .
		  ?country rdfs:label ?label .
		}`, args.Country)

	// 2. Prepare the Request
	apiURL := "https://query.wikidata.org/sparql"
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return GetCountryResult{}, err
	}

	q := req.URL.Query()
	q.Add("query", sparqlQuery)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()
	
	// Wikidata requires a User-Agent header
	req.Header.Set("User-Agent", "EtymologyScorerBot/1.0 (contact: kostas@loupax.com)")

	// 3. Execute the Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetCountryResult{}, err
	}
	defer resp.Body.Close()

	// 4. Parse the Wikidata Response
	var result struct {
		Results struct {
			Bindings []struct {
				Label struct {
					Value string `json:"value"`
				} `json:"label"`
			} `json:"bindings"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GetCountryResult{}, err
	}
uniqueMap := make(map[string]bool)
    
    // Create a transformer to strip accents (Mn = Mark, Nonspacing)
    t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

    for _, b := range result.Results.Bindings {
        // 1. Lowercase
        raw := strings.ToLower(b.Label.Value)
        
        // 2. Remove Accents/Diacritics
        clean, _, _ := transform.String(t, raw)
        
        // 3. Remove non-letters (optional, helps with "Greece!" vs "Greece")
        clean = strings.Map(func(r rune) rune {
            if unicode.IsLetter(r) {
                return r
            }
            return -1
        }, clean)

        if clean != "" {
            uniqueMap[clean] = true
        }
    }

    var names []string
    for name := range uniqueMap {
        names = append(names, name)
    }

		fmt.Println("names fetched:", names)
    return GetCountryResult{Names: names}, nil
}
