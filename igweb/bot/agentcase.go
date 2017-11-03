package bot

import (
	"math/rand"
	"time"

	"github.com/james-bowman/nlp"
	"gonum.org/v1/gonum/mat"
)

func randomNumber(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

type AgentCase struct {
	Bot
	name            string
	title           string
	thumbnailPath   string
	knowledgeBase   map[string]string
	knowledgeCorpus []string
	sampleQuestions []string
}

func NewAgentCase() *AgentCase {
	agentCase := &AgentCase{name: "Case", title: "Resident Isomorphic Gopher Agent", thumbnailPath: "/static/images/chat/Case.png"}
	agentCase.initializeIntelligence()
	return agentCase
}

func (a *AgentCase) Name() string {
	return a.name
}

func (a *AgentCase) Title() string {
	return a.title
}

func (a *AgentCase) ThumbnailPath() string {
	return a.thumbnailPath
}

func (a *AgentCase) SetName(name string) {
	a.name = name
}

func (a *AgentCase) SetTitle(title string) {
	a.title = title
}

func (a *AgentCase) SetThumbnailPath(thumbnailPath string) {
	a.thumbnailPath = thumbnailPath
}

func (a *AgentCase) Greeting() string {

	sampleQuestionIndex := randomNumber(0, len(a.sampleQuestions))
	greeting := "Hi there! I'm Case. You can ask me a question on Isomorphic Go. Such as...\"" + a.sampleQuestions[sampleQuestionIndex] + "\""
	return greeting

}

func (a *AgentCase) Reply(query string) string {

	var result string

	vectoriser := nlp.NewCountVectoriser(true)
	transformer := nlp.NewTfidfTransformer()

	// set k (the number of dimensions following truncation) to 4
	reducer := nlp.NewTruncatedSVD(4)

	// Transform the corpus into an LSI fitting the model to the documents in the process
	matrix, _ := vectoriser.FitTransform(a.knowledgeCorpus...)
	matrix, _ = transformer.FitTransform(matrix)
	lsi, _ := reducer.FitTransform(matrix)

	// run the query through the same pipeline that was fitted to the corpus and
	// to project it into the same dimensional space
	matrix, _ = vectoriser.Transform(query)
	matrix, _ = transformer.Transform(matrix)
	queryVector, _ := reducer.Transform(matrix)

	// iterate over document feature vectors (columns) in the LSI and compare with the
	// query vector for similarity.  Similarity is determined by the difference between
	// the angles of the vectors known as the cosine similarity
	highestSimilarity := -1.0
	var matched int
	_, docs := lsi.Dims()
	for i := 0; i < docs; i++ {
		similarity := nlp.CosineSimilarity(queryVector.(mat.ColViewer).ColView(0), lsi.(mat.ColViewer).ColView(i))
		if similarity > highestSimilarity {
			matched = i
			highestSimilarity = similarity
		}
	}

	if highestSimilarity == -1 {
		result = "I don't know the answer to that one."
	} else {
		result = a.knowledgeBase[a.knowledgeCorpus[matched]]
	}

	if result == "" {
		result = "Sorry, I don't know what you're talking about. My knowledge is limited."
	}

	return result

}

func (a *AgentCase) initializeIntelligence() {

	a.knowledgeBase = map[string]string{
		"isomorphic go isomorphic go web applications":                                            "Isomorphic Go is the methodology to create isomorphic web applications using the Go (Golang) programming language. An isomorphic web application, is a web application, that contains code which can run, on both the web client and the web server.",
		"kick recompile code restart web server instance instant kickstart lightweight mechanism": "Kick is a lightweight mechanism to provide an instant kickstart to a Go web server instance, upon the modification of a Go source file within a particular project directory (including any subdirectories). An instant kickstart consists of a recompilation of the Go code and a restart of the web server instance. Kick comes with the ability to take both the go and gopherjs commands into consideration when performing the instant kickstart. This makes it a really handy tool for isomorphic golang projects.",
		"starter code starter kit":                                                                "The isogoapp, is a basic, barebones web app, intended to be used as a starting point for developing an Isomorphic Go application. Here's the link to the github page: https://github.com/isomorphicgo/isogoapp",
		"idiot stupid dumb dummy don't know anything":                                             "Please don't question my intelligence, it's artificial after all!",
		"more info a topic where can i find talk lecture presentation":                            "Watch the Isomorphic Go talk by Kamesh Balasubramanian at GopherCon India: https://youtu.be/zrsuxZEoTcs",
		"benefits of the technology significance of the technology importance of the technology":  "Here are some benefits of Isomorphic Go: Unlike JavaScript, Go provides type safety, allowing us to find and eliminate many bugs at compile time itself. Eliminates mental context-shifts between back-end and front-end coding. Page loading prompts are not necessary.",
		"perform routing web app register routes define routes":                                   "You can implement client-side routing in your web application using the IsoKit Router preventing the dreaded full page reload.",
		"render templates perform template rendering":                                             "Use template sets, a set of project templates that are persisted in memory and are available on both the server-side and the client-side",
		"cogs reusable components react-like react":                                               "Cogs are reuseable components in an Isomorphic Go web application.",
	}

	a.knowledgeCorpus = make([]string, 1)

	for k, _ := range a.knowledgeBase {
		a.knowledgeCorpus = append(a.knowledgeCorpus, k)
	}

	a.sampleQuestions = []string{"What is isomorphic go?", "What are the benefits of this technology?", "Does isomorphic go offer anything react-like?", "How can I recompile code instantly?", "How can I perform routing in my web app?", "Where can I get starter code?", "Where can I watch a talk on this subject?"}

}