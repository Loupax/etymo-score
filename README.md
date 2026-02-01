# 🌍 EtymoScore: Country Etymology Agent

**EtymoScore** is a specialized AI agent built with the **Google Agent Development Kit (ADK)** in Go. It calculates a "Linguistic Variation Score" for any country by analyzing historical names, exonyms, and endonyms to determine their unique etymological roots (cognates).

## 🚀 Key Features

- **Cognate-Based Grouping:** Instead of just counting names, the agent uses Gemini's reasoning to group variations by their common ancestor (e.g., *Spain* and *España* are grouped as one; *Hungary* and *Magyarország* are separated).
- **Gemini 3 Pro Reasoning:** Leverages the latest "Thinking" capabilities to handle complex historical linguistics.
- **Google Search Grounding:** Uses real-time search to find rare, historical, or regional variations of country names.
- **Structured Output:** Delivers a clean, programmatic JSON response for easy integration into other Go services.

## 🛠 Tech Stack

- **Language:** Go (Golang) 1.22+
- **Framework:** [Google ADK](https://google.golang.org/adk)
- **Model:** Gemini 3 Pro Preview
- **Tools:** Google Search (Grounding)

## 📋 Prerequisites

- Go installed on your system.
- A [Google AI Studio](https://aistudio.google.com/) API Key.

## ⚙️ Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/etymo-score.git
   cd etymo-score
   ```

2. **Set your API Key:**
   ```bash
   export GOOGLE_API_KEY="your_actual_api_key_here"
   ```

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

## 🏃 Usage

Run the agent via the CLI launcher provided by the ADK:

```bash
go run main.go "What is the variation score for Greece?"
```

### Example Programmatic Result
The agent provides a response structured like this:
```json
{
  "score": 2,
  "explanation": "Group 1: Greece, Griechenland (Latin root Graecia). Group 2: Ελλάδα, Hellas (Hellenic root)."
}
```

## 🧠 How it Works: The Cognate Logic

The agent is instructed to follow a rigid linguistic rule-set:

1. **Search:** Fetch variations (e.g., Germany -> Germany, Deutschland, Allemagne, Németország).
2. **Analyze Roots:**
   - *Germanic* (Latin) -> Germany
   - *Alemanni* (Tribal) -> Allemagne
   - *Slavic* (Mute) -> Németország
3. **Score:** The agent identifies **3 unique roots**, resulting in a score of **3**.

## 📄 License
This project is licensed under the MIT License.
