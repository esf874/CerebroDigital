# Proposal

Create a solution for a note service enhanced with AI. The app will use a text-based format (markdown) which integrates well with LLMs.

# Architecture

## Backend

We propose to build the backend on the Go language, which is highly used for building cloud services and has a very small ramp. The student should have either familiarity with the language or course the [Tour of Go](https://go.dev/tour/welcome/1). 

The service must implement CRUD (Create, Retrieve, Update and Delete) notes methods, and we recommend a plain, tag hierarchy for the notes, which is common for apps such as [Obsidian](https://obsidian.md/) and [Logseq](https://logseq.com/). The user creates notes, all with the same hierarchy level, and instead build relations between the notes using references and tags. Leaving the information to structure itself, instead of the more traditional directory tree. This will be usefull for gathering relevant information for the LLM.

We recommend to store this information directly on a mongo database for simplicity. 

![](./obsidian-graph-view-notes-2048x1111.png)

When the user adds or updates a note, the backend must analyze its content looking for references to other notes and tags. This is an extension over markdown and must be defined by the application, for example, Logseq uses __double square brakets__ for references **[[this is a reference]]** and sharp+keyword for **#tags**. 

The backend should store this information on a doubly connected edge list graph structure. When the user asks the model, the backend will gather the information with related tags and references and pass it down the model for [retrieval augmented generation](https://aws.amazon.com/what-is/retrieval-augmented-generation/). 

## LLM

For the LLM, we recommend adhering to the OpenAI api, which is becoming an standard. Many models which can be run locally now offer this api, for example [ollama](https://docs.ollama.com/api/openai-compatibility#openai-compatibility) or the [VLLM framework](https://docs.vllm.ai/en/latest/getting_started/quickstart.html#openai-compatible-server).

## Frontend

This should be an optional and complementary part for the project, since a lot of time could be lost here and the rest of the project could be tested directly calling the API. If necessary, there are libraries already in JS for [markdown rendering](https://marked.js.org/demo/?text=Marked%20-%20Markdown%20Parser%0A%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%3D%0A%0A%5BMarked%5D%20lets%20you%20convert%20%5BMarkdown%5D%20into%20HTML.%20%20Markdown%20is%20a%20simple%20text%20format%20whose%20goal%20is%20to%20be%20very%20easy%20to%20read%20and%20write%2C%20even%20when%20not%20converted%20to%20HTML.%20%20This%20demo%20page%20will%20let%20you%20type%20anything%20you%20like%20and%20see%20how%20it%20gets%20converted.%20%20Live.%20%20No%20more%20waiting%20around.%0A%0AHow%20To%20Use%20The%20Demo%0A-------------------%0A%0A1.%20Type%20in%20stuff%20on%20the%20left.%0A2.%20See%20the%20live%20updates%20on%20the%20right.%0A%0AThat%27s%20it.%20%20Pretty%20simple.%20%20There%27s%20also%20a%20drop-down%20option%20above%20to%20switch%20between%20various%20views%3A%0A%0A-%20**Preview%3A**%20%20A%20live%20display%20of%20the%20generated%20HTML%20as%20it%20would%20render%20in%20a%20browser.%0A-%20**HTML%20Source%3A**%20%20The%20generated%20HTML%20before%20your%20browser%20makes%20it%20pretty.%0A-%20**Lexer%20Data%3A**%20%20What%20%5Bmarked%5D%20uses%20internally%2C%20in%20case%20you%20like%20gory%20stuff%20like%20this.%0A-%20**Quick%20Reference%3A**%20%20A%20brief%20run-down%20of%20how%20to%20format%20things%20using%20markdown.%0A%0AWhy%20Markdown%3F%0A-------------%0A%0AIt%27s%20easy.%20%20It%27s%20not%20overly%20bloated%2C%20unlike%20HTML.%20%20Also%2C%20as%20the%20creator%20of%20%5Bmarkdown%5D%20says%2C%0A%0A%3E%20The%20overriding%20design%20goal%20for%20Markdown%27s%0A%3E%20formatting%20syntax%20is%20to%20make%20it%20as%20readable%0A%3E%20as%20possible.%20The%20idea%20is%20that%20a%0A%3E%20Markdown-formatted%20document%20should%20be%0A%3E%20publishable%20as-is%2C%20as%20plain%20text%2C%20without%0A%3E%20looking%20like%20it%27s%20been%20marked%20up%20with%20tags%0A%3E%20or%20formatting%20instructions.%0A%0AReady%20to%20start%20writing%3F%20%20Either%20start%20changing%20stuff%20on%20the%20left%20or%0A%5Bclear%20everything%5D(%2Fdemo%2F%3Ftext%3D)%20with%20a%20simple%20click.%0A%0A%5BMarked%5D%3A%20https%3A%2F%2Fgithub.com%2Fmarkedjs%2Fmarked%2F%0A%5BMarkdown%5D%3A%20http%3A%2F%2Fdaringfireball.net%2Fprojects%2Fmarkdown%2F%0A&options=%7B%0A%20%22async%22%3A%20false%2C%0A%20%22breaks%22%3A%20false%2C%0A%20%22extensions%22%3A%20null%2C%0A%20%22gfm%22%3A%20true%2C%0A%20%22hooks%22%3A%20null%2C%0A%20%22pedantic%22%3A%20false%2C%0A%20%22silent%22%3A%20false%2C%0A%20%22tokenizer%22%3A%20null%2C%0A%20%22walkTokens%22%3A%20null%0A%7D&version=16.4.1). 

## Environment

We recommend using [Nix](https://nixos.org/download/) for development environments. We provide already a shell.nix file, which can be used with the command `nix-shell` or with [direnv](https://direnv.net/) automatically when cd'ing. With the second and the recommended extensions, vscode should automatically access the environment. 

Install direnv by following the instructions for the linux distribution of your choice and remember to [hook it into your shell](https://direnv.net/docs/hook.html). 

Install the recommended extensions by searching `@recommended` on the extensions pannel. 

You can launch the server with `cd backend && air .`

