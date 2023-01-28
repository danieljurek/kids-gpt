## Kids GPT 🔰🤖

A console application for macOS that uses dictation and `say` which gives a
child the ability to play with GPT-3. Use configuration `.yml` files.

![Screeshot of Kids GPT](https://user-images.githubusercontent.com/2158838/215277967-7335deb6-a91e-4fb5-9597-e2688f414361.png)

### Pre-requisites

* An OpenAI account and API key
* Go tools


### How to use

Set the API key

```bash
export OPENAI_KEY=<your key>
```

Build

```bash
go build
```

Run
```
./kids-gpt story.yml
```

GPT-3 will greet you according to its prompt and you will have an opportunity to
respond by speaking or typing.

During prompts press the microphone key (F5 on some macs) to launch dictation.
You may have to enable dictation on your device.

Press `return` to submit your answer to GPT-3 and hear a response.

### Configuration

See the yaml files for configuration.

### ⚠️ Help wanted!

This tool is made primarily for me and my daughter to learn how to interact with
GPT-3. You should try it with your kids, too!

What I've observed so far:

* macOS dictation is good with adults and not great at understanding kids. I'm
interested in better dictation for kids 4-5 years old.
* Your child needs to speak clearly into the microphone (facing the monitor
generally works well). This is best accomplished in a quiet room with few
distractions.
* Keep an eye on what the computer is saying to your kid and intervene in the
dictation process if things get out of hand.
