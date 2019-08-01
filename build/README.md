Build
=====

An image used for packaging and distributing MySQL database images.

## Example

```bash
$ export IMAGE=docker.io/my/image

$ docker run -it -v $HOME/.docker:/kaniko/.docker \
                 -v $(pwd)/example:/workspace \
                 skpr/mtk-build:latest --context=/workspace \
                                       --dockerfile=/Dockerfile \
                                       --single-snapshot \
                                       --destination=$IMAGE:latest \
                                       --destination=$IMAGE:$(date +%F) \
                                       --verbosity fatal
```

Note: See `build.dockerfiler` for what can be configured with `--build-arg`
