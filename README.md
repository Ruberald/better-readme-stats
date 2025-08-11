# 📊 GitHub Top Languages by Commit Count

This project generates a **Top Languages SVG card** based on your **commit history**, not code size.  
It queries the GitHub GraphQL API to fetch your repositories, counts commits per primary language, and produces a visual card you can embed in your README or website.

---

## 🚀 Features

- Fetches your GitHub repositories via GraphQL API.
- Calculates **language percentages based on commit count**.
- Displays the **top 7 languages** plus an `"Others"` category.
- SVG output is styled with GitHub language colors.
- Lightweight Go HTTP server for serving the SVG dynamically.

---

## 📥 Installation (Local Development)

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/top-langs-commits.git
   cd top-langs-commits
    ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Create a `colors.json` file containing GitHub language colors (from the [linguist repository](https://github.com/github/linguist)).

4. Set up your GitHub personal access token (with `repo` and `read:user` scopes):

   ```bash
   export GITHUB_TOKEN=your_token_here
   ```

5. Run the server:

   ```bash
   go run main.go
   ```

---

## 🌐 Live Hosted Version

You can see it live without setting anything up!

Just use this URL:

```
https://better-readme-stats-ruberald8800-b77rzyn5.leapcell.dev/stats?username=YOUR_GITHUB_USERNAME
```

### Example with my username:

![Top Languages by Commits](https://better-readme-stats-ruberald8800-b77rzyn5.leapcell.dev/stats?username=Ruberald)

To embed in your own README, just add:

```markdown
![Top Languages by Commits](https://better-readme-stats-ruberald8800-b77rzyn5.leapcell.dev/stats?username=YOUR_GITHUB_USERNAME)
```

---

## 📊 How Percentages Are Calculated

### **This Project**

* Percentages are **commit-based**:

  ```
  percent = (commits for this language / total commits) * 100
  ```
* Example: If you have 400 commits in Go and 100 in Python (total 500), Go = 80%, Python = 20%.
* Top 7 languages shown, rest grouped into `"Others"`.

### **[anuraghazra/github-readme-stats](https://github.com/anuraghazra/github-readme-stats)**

* Percentages are **size + repo count based**:

  ```
  ranking_index = (byte_count ^ size_weight) * (repo_count ^ count_weight)
  ```
* Favors languages with:

  * Large total code size in your repos.
  * More repositories using that language.

### **Key Difference**

* **This project** → Measures *how often you work with a language* (**activity**).
* **Readme Stats** → Measures *how much of your codebase is in that language* (**size & popularity**).

---

## 🛠 Tech Stack

* **Go** — server, API fetching, SVG templating
* **GitHub GraphQL API** — repository and commit data
* **HTML Templates** — SVG generation
* **JSON** — language color mapping

---

## 📄 License

MIT License — feel free to modify and use in your own projects.

```

If you want, I can also add a **"How It Works"** diagram to visually show the API call → commit counting → SVG generation flow, which would make this README even more eye-catching.
```
