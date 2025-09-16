# AI Safety

This document provides an overview of LLM-based guardrails and explores a few prominent solutions. 

## LLM Guardrails

Guardrails for LLMs provide input and output moderation to detect and mitigate a wide range of potential harms. These guardrails can encompass various techniques, including content filtering, prompt engineering, and output validation, all aimed at aligning LLM behavior with desired safety guidelines and legal compliance. 

A few leading open-source and closed-source solutions include OpenAI content moderation API, Perspective API, and Llama Guard. These solutions primarily rely on another pre-trained LLM to understand the input and output of the target LLM contextually.

## Available Guardrails Solutions

### 1. OpenAI Content Moderation API

OpenAI's moderation API, accessible via their "moderations" endpoint, offers free access to GPT-based classifiers. These classifiers are designed to detect a broad set of categories of harmful content. At the time of writing, OpenAI supports the following categories: 

```
harassment, harassment/threatening, hate, hate/threatening, illicit, illicit/violent, self-harm, self-harm/intent, self-harm/instructions, sexual, sexual/minors, violence, violence/graphic.
```

Currently, the “moderations” endpoint can be used with two models whereas the newest model i.e. omni supports multi-modal inputs and more categorization options. The list of categories and model support is available at [content-classifications](https://platform.openai.com/docs/guides/moderation/overview#content-classifications). 

### 2. Llama Guard

Llama Guard is an LLM-based input-output safeguard model, fine-tuned on data labeled according to their taxonomy. Llama Guard includes the applicable taxonomy as the input and uses instruction tasks for classification. This enables users to customize the model input, facilitating adaption to other taxonomies appropriate for their use case with zero-shot or few-shot prompting. Furthermore, Llama Guard uses a slightly different set of categories of harmful content than Open AI’s moderation API.
