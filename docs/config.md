# Configuration
This project uses environment variables to get its config. You can find a sample config in [config/example.env](../config/example.env)

You can either set the environment variables globally via your shell environment, or you can set them in a file, and
then set that file as the environment before runnning the program.

## Deployment
For deployments that use SystemD you can simply set your prod env file as the environment and it will work as intended.
Any technique that loads the env file should also work.

## Development
For development, there are multiple ways to run this project with a specific environment.

### IDE
For IntelliJ IDEs, you can install the `EnvFile` plugin, and then in your run config, you should see a tab that will
allow you to set an env file for that specific config.

There is probably something like this for other editors but I do not use them.

### Command line
To run the project with the variables loaded, you can either just put all the variables in your
`.zshrc`/`.bashrc`/`.whateverrc`to load them globally. I do not recommend this approach because the variables will be
available globally and you may have conflicts with other programs.

What I would recommend is to use a simple bash trick that will let you load the environment, in a way that isolates it
to the running program.
```bash
(set -a source example.env set +a ./passkey-server)
```
Make sure to replace `example.env` with your own config.

