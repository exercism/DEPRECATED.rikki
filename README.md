# Rikki

Rikki is the friendly neighborhood robot code reviewer.

Rikki's comments are open source, and we welcome improvements of all kinds.

If you find typos, or if something is poorly worded, confusing, or incorrect,
please file an issue or submit a pull request.

## Dependencies

* Ruby Analyzer
    * **github repo**: https://github.com/exercism/rikki-ruby-analyzer
    * **deployed to**: heroku (https://git.heroku.com/exercism-ruby-analyzer.git)
    * **URL**: http://ruby-analyzer.exercism.io
* Crystal Analyzer
    * **github repo**: https://github.com/mhelmetag/crystal-analyzer
    * **deployed to**: heroku (https://git.heroku.com/exercism-crystal-analyzer.git)
    * **URL**: http://crystal-analyzer.exercism.io

The Ruby analyzer uses analyzers/rules defined in the [exercism-analysis repository](https://github.com/JacobNinja/exercism-analysis),
and responds with a list of violations. Each violation consists of type and a
possible list of keys.

```json
{
  "results":[
    {
      "type": "enumerable_condition",
      "keys": [
        "enumerable_condition"
      ]
    }
  ]
}
```

## How it works

When someone submits a solution to exercism, the uuid of the submission gets
placed in a queue implemented in redis.

Rikki pulls the uuid off the queue, makes a request to exercism.io's API to get the
code, then either processes it locally, or submits the code to one of the analyzer APIs,
depending on the language.

* Go: analyzes locally
* Crystal: submits it to the crystal analyzer
* Ruby: submits it to the ruby analyzer

The `comments/` directory of the rikki project contains a directory for each
`type`, and a markdown file for each `key`.

The worker chooses one key at random and submits the contents of the markdown
file as a comment to the exercism.io.

## Usage

```bash
$ go install github.com/exercism/rikki
$ rikki
```

By default the worker runs against

- **redis**: localhost:6379
- **exercism**: localhost:4567
- **crystal-analyzer**: localhost:3000
- **ruby-analyzer**: localhost:8989

These can be overridden using command-line flags:

```bash
$ ./rikki \
    -redis=redis://user:pass@host:port/db/ \
    -exercism=http://exercism.io \
    -crystal-analyzer=http://crystal-analyzer.exercism.io \
    -ruby-analyzer=http://ruby-analyzer.exercism.io
```

## Enqueuing a Job

Start the console in exercism, find the uuid (a.k.a. `key`)  of the submission
you wish to process, and enqueue the job with:

```
Jobs::Analyze.perform_async(uuid)
```
