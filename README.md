# A go wrapper for OpenAI Platform APIs

This is a wrapper around OpenAI APIs that means to simplify the creation of
tooling that interacts with OpenAI and other Go code in fun, novel ways. There
is also a generic fun cli in `cmd/openai` that exposes library code to the CLI
as well as providing some usage examples for the library.

If you need a token, head over to https://platform.openai.com/account/api-keys,
and sign up if you need to. If you haven't signed up already, you usually get
$20ish in credit, which should get you pretty far.


[![Go Reference](https://pkg.go.dev/badge/github.com/andrewstuart/openai.svg)](https://pkg.go.dev/github.com/andrewstuart/openai)
![Go](https://github.com/andrewstuart/openai/actions/workflows/go.yml/badge.svg)


You may also notice that there are plenty of pointers for optional fields in
Request structs, but none of the common helper funcs to return pointers. I
created one module to rule them all (there may be others, but it's 3 lines and I
honestly didn't want to spend the time looking). I suggest
[github.com/andrewstuart/p](github.com/andrewstuart/p) for your "pointer to anything, even a literal or
const" needs.

```bash
go install github.com/andrewstuart/openai/cmd/openai@latest

echo "token: $MY_TOKEN" >> $HOME/.config/openai.yaml
# Or export TOKEN in your environment variables.
# add org: <your org id> if you are a member of multiple organizations.

openai chat --personality "Lady Gaga"
openai chat --prompt "You answer in the style of a super chill surfer from southern california."
openai image "Marvin the Paranoid Android giving a speech." -f marvin.jpg
openai variations -n 5 marvin-001.jpg
openai complete 'A description of a painting of a perfect day' | openai image -f self.jpg -
```

For the best current go examples, see the [CLI files](cmd/openai/cmd). 

Current support:

- [x] Chat
- [x] Completion (code, text, etc)
- [x] Edit
- [x] Images
- [x] Image variations
- [x] Audio transcription
- [ ] Upload/manage Files (for fine tuning)
- [ ] Fine tune models
- [ ] Moderations (check for OpenAI content policy violations)
- [ ] More soon
