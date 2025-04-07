---
slug: /quick-start
description: Civic's quick start guide to help you get started with your CV.
keyword: cv, resume, curriculum vitae, civic, job, work
---

# Quick Start

Before jumping to all the jargon, let's just go through a simple set of
instructions to generate the first draft of our CV. This step helps us
understand if Civic is really what it claims to be. If you didn't enjoy the
process, just move on, no hard feelings 🥲.

## Step 1: Install Civic

Make sure you have installed Civic on your machine:

```bash
civic --version
```

If you see the version number, you are good to go. If not, please follow the
[installation guide](./02_installation) to install Civic on your machine.

## Step 2: Create a Configuration File

You can use the following command to create a sample configuration file:

```bash
civic schema init
```

This command create a new `.civic.yaml` file in your current working
directory. The file contains a sample configuration file that has arbitrary
data for each section of Civic's schema. You can view the file and change
whatever information you like.

However, the point of this guide is to make you familiar with the steps.

Note that for the sake of simplicity, this sample configuration uses the
simplest template file possible. So, don't be disappointed about your end
result. You can brush it up later.

## Step 3: Generate Your CV

Now, you can generate your CV with the following command:

```bash
civic generate
```

This command can be configured to work with a specific configuration file or
generate a different output. However, since we have a `.civic.yaml` file in
the working directory, we can use the `generate` command without any arguments.

Once the execution is completed, you should be able to see a file called
`civic.pdf` in the working directory. The exact name and path of the file
will be mentioned in the command output as well.

## Step 4: Review Your CV

Open the `civic.pdf` file and review your CV. If you are happy with the
result, you can share it with others. If not, there are two ways. Do you
have any feedback that could improve the process and make you happy? If yes,
please share it with us. If not, then this is farewell 👋.

## Step 5: Sharpen your CV

Now that you're happy and ready to move forward 🎉, navigate through the rest
of the documentation to learn more about the schema, CLI commands, and how
to customize your templates.
