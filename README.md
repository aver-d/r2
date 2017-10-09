# r2

HTTP router. Provided here for demonstration/educational purposes only.

During the course of every person's life, there's an outside chance that for some obscure reason he or she may want to implement an HTTP router rather than use an existing robust alternative.

Some code here does that.

This is written in Go, but probably best suited to be ported to some other language.
The router does not implement the `http.Handler` interface.
Instead handlers can be written in a compact way like below.
Handlers are passed an environment `Env`, which contains the request, response and any other stuff you might want.

    package main

    import (
        "fmt"
        "github.com/aver-d/r2"
        "net/http"
    )

    func main() {
        router := r2.NewRouter("/api")

        router.Get("/hello/:name", func(env *r2.Env) {
            fmt.Fprint(env.W, "Hello "+env.Path["name"])
        })

        http.ListenAndServe("localhost:4444", router)
    }

Variables in paths can be accessed as-is or matched against a regular expression with a special `!` syntax following the declaration.
For example, `:name![dD].+` will match only names begining with d.
Special values `int` and `float` are provided to match numbers as a more descriptive alternative to defining a regular expression.
So, `/:age!int` will match only if the value of `age` can be converted to an integer.

I mention above that a sensible idea is normally to use an existing, battle-tested router.
A commonly-used choice is one by [Julien Schmidt], and
that router at one point used the Github api as test data.

Now, returning to my router here, we can also use that Github api data as an example.
With r2 a recursive representation of the router can be sent to stdout using `router.Print()`.

So, here is the Github api (as of 2015) when used with r2.
Each endpoint is followed by a list of `(METHOD, function_name)` pairs.
In this example, the function name is always shown using a fake handler called `f`
(I, of course, couldn't know what the actual Github handlers are).


    /
    ├───applications
    │    └───:client_id
    │        └───tokens DELETE f
    │            └───:access_token GET f DELETE f
    ├───authorizations GET f POST f
    │    ├───:id GET f PATCH f DELETE f
    │    └───clients
    │        └───:client_id PUT f
    ├───emojis GET f
    ├───events GET f
    ├───feeds GET f
    ├───gists GET f POST f
    │    ├───:id GET f PATCH f DELETE f
    │    │    ├───forks POST f
    │    │    └───star PUT f DELETE f GET f
    │    ├───public GET f
    │    └───starred GET f
    ├───gitignore
    │    └───templates GET f
    │        └───:name GET f
    ├───issues GET f
    ├───legacy
    │    ├───issues
    │    │    └───search
    │    │        └───:owner
    │    │            └───:repository
    │    │                └───:state
    │    │                    └───:keyword GET f
    │    ├───repos
    │    │    └───search
    │    │        └───:keyword GET f
    │    └───user
    │        ├───email
    │        │    └───:email GET f
    │        └───search
    │            └───:keyword GET f
    ├───markdown POST f
    │    └───raw POST f
    ├───meta GET f
    ├───networks
    │    └───:owner
    │        └───:repo
    │            └───events GET f
    ├───notifications GET f PUT f
    │    └───threads
    │        └───:id GET f PATCH f
    │            └───subscription GET f PUT f DELETE f
    ├───orgs
    │    └───:org GET f PATCH f
    │        ├───events GET f
    │        ├───issues GET f
    │        ├───members GET f
    │        │    └───:user GET f DELETE f
    │        ├───public_members GET f
    │        │    └───:user GET f PUT f DELETE f
    │        ├───repos GET f POST f
    │        └───teams GET f POST f
    ├───rate_limit GET f
    ├───repos
    │    └───:owner
    │        └───:repo GET f PATCH f DELETE f
    │            ├───:archive_format
    │            │    └───:ref GET f
    │            ├───assignees GET f
    │            │    └───:assignee GET f
    │            ├───branches GET f
    │            │    └───:branch GET f
    │            ├───collaborators GET f
    │            │    └───:user GET f PUT f DELETE f
    │            ├───comments GET f
    │            │    └───:id GET f PATCH f DELETE f
    │            ├───commits GET f
    │            │    └───:sha GET f
    │            │        └───comments GET f POST f
    │            ├───contents
    │            │    └───*path DELETE f GET f PUT f
    │            ├───contributors GET f
    │            ├───downloads GET f
    │            │    └───:id GET f DELETE f
    │            ├───events GET f
    │            ├───forks GET f POST f
    │            ├───git
    │            │    ├───blobs POST f
    │            │    │    └───:sha GET f
    │            │    ├───commits POST f
    │            │    │    └───:sha GET f
    │            │    ├───refs GET f POST f
    │            │    │    └───*ref GET f PATCH f DELETE f
    │            │    ├───tags POST f
    │            │    │    └───:sha GET f
    │            │    └───trees POST f
    │            │        └───:sha GET f
    │            ├───hooks GET f POST f
    │            │    └───:id GET f PATCH f DELETE f
    │            │        └───tests POST f
    │            ├───issues GET f POST f
    │            │    ├───:number GET f PATCH f
    │            │    │    ├───comments GET f POST f
    │            │    │    ├───events GET f
    │            │    │    └───labels DELETE f GET f POST f PUT f
    │            │    │        └───:name DELETE f
    │            │    ├───comments GET f
    │            │    │    └───:id GET f PATCH f DELETE f
    │            │    └───events GET f
    │            │        └───:id GET f
    │            ├───keys GET f POST f
    │            │    └───:id GET f PATCH f DELETE f
    │            ├───labels GET f POST f
    │            │    └───:name PATCH f DELETE f GET f
    │            ├───languages GET f
    │            ├───merges POST f
    │            ├───milestones POST f GET f
    │            │    └───:number DELETE f GET f PATCH f
    │            │        └───labels GET f
    │            ├───notifications GET f PUT f
    │            ├───pulls GET f POST f
    │            │    ├───:number GET f PATCH f
    │            │    │    ├───comments GET f PUT f
    │            │    │    ├───commits GET f
    │            │    │    ├───files GET f
    │            │    │    └───merge GET f PUT f
    │            │    └───comments GET f
    │            │        └───:number GET f PATCH f DELETE f
    │            ├───readme GET f
    │            ├───releases GET f POST f
    │            │    └───:id GET f PATCH f DELETE f
    │            │        └───assets GET f
    │            ├───stargazers GET f
    │            ├───stats
    │            │    ├───code_frequency GET f
    │            │    ├───commit_activity GET f
    │            │    ├───contributors GET f
    │            │    ├───participation GET f
    │            │    └───punch_card GET f
    │            ├───statuses
    │            │    └───:ref GET f POST f
    │            ├───subscribers GET f
    │            ├───subscription GET f PUT f DELETE f
    │            ├───tags GET f
    │            └───teams GET f
    ├───repositories GET f
    ├───search
    │    ├───code GET f
    │    ├───issues GET f
    │    ├───repositories GET f
    │    └───users GET f
    ├───teams
    │    └───:id GET f PATCH f DELETE f
    │        ├───members GET f
    │        │    └───:user GET f PUT f DELETE f
    │        └───repos GET f
    │            └───:owner
    │                └───:repo GET f PUT f DELETE f
    ├───user GET f PATCH f
    │    ├───emails GET f POST f DELETE f
    │    ├───followers GET f
    │    ├───following GET f
    │    │    └───:user GET f PUT f DELETE f
    │    ├───issues GET f
    │    ├───keys GET f POST f
    │    │    └───:id GET f PATCH f DELETE f
    │    ├───orgs GET f
    │    ├───repos GET f POST f
    │    ├───starred GET f
    │    │    └───:owner
    │    │        └───:repo PUT f DELETE f GET f
    │    ├───subscriptions GET f
    │    │    └───:owner
    │    │        └───:repo PUT f DELETE f GET f
    │    └───teams GET f
    └───users GET f
        └───:user GET f
            ├───events GET f
            │    ├───orgs
            │    │    └───:org GET f
            │    └───public GET f
            ├───followers GET f
            ├───following GET f
            │    └───:target_user GET f
            ├───gists GET f
            ├───keys GET f
            ├───orgs GET f
            ├───received_events GET f
            │    └───public GET f
            ├───repos GET f
            ├───starred GET f
            └───subscriptions GET f

License: MIT

[Julien Schmidt]: https://github.com/julienschmidt/httprouter

