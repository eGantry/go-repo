# go-repo
Go game in Go Language

This started out as an innocent chat with ChatGPT about the origins of both the Go board game, and the Go programming language,
then it turned into a great excuse to gain a little knowledge about both, by merging them into one project.

Each working phase of the chat shows up here as a commit.  Unlike a lot of my chat-assisted Java projects, the code recommendations
from this chat worked without a lot of tweaking needed.

If you're following along with the chat, you'll notice I bailed on the recommended fancy text option, because wrestling with the 
issues it created was distracting from the keeping the main functionality working.

If this were a "real" project and fancy text were a spec, I would capture this as a known issue, with simpler text as a workaround.

With the latest commit, I now have a nice-looking, working, playable Go board for two players.

If you wish to implement this board yourself, you can use this command in the project folder to run it:

go run main.go ko_rule.go
