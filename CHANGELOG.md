# Changelog

## 1.0.0 (2026-04-09)


### Features

* add application metrics with expvar integration and request tracking ([db71844](https://github.com/ncondes/go-social/commit/db7184415ac703cb8997b36c200198ac09e54060))
* add authorization layer with role-based access control for posts and comments ([128f709](https://github.com/ncondes/go-social/commit/128f7099741d8692fe65df894ee2fa4e1e6eea9a))
* add cors support, build frontend activation ui, and enhance development setup ([110186b](https://github.com/ncondes/go-social/commit/110186b8d26ff65ce0cc586b84e2c820258fed0e))
* add logging interface and integrate zap logger into handlers ([3d484db](https://github.com/ncondes/go-social/commit/3d484dbe47bdf3676368de5c2be357691cc39ca2))
* add redis caching layer and enhance makefile ([c2a6e1f](https://github.com/ncondes/go-social/commit/c2a6e1f49cb10b98a33c81bdc777ed69672290aa))
* add release-please workflow and enhance ci pipeline with swagger and migration steps ([62a43a5](https://github.com/ncondes/go-social/commit/62a43a570057c0b8085d38e41e804652e79d16df))
* database setup and clean architecture bases ([86a6b46](https://github.com/ncondes/go-social/commit/86a6b462a0e34996bfb320054a1e66edfe75ddb7))
* implement jwt authentication with token generation and protected routes ([28c61e1](https://github.com/ncondes/go-social/commit/28c61e160be4ca7b79ada51a7d56e7f36e1d377a))
* implement rate limiting with graceful shutdown ([833c22c](https://github.com/ncondes/go-social/commit/833c22cf5214cb1de7c89a300e189dda825f33e9))
* implement user registration and activation with invitation tokens ([b8900cb](https://github.com/ncondes/go-social/commit/b8900cb55d1885ff678f03ef3fa9d82c59490e12))
* integrate sendgrid for email invitations and add user deletion functionality ([9bf02fc](https://github.com/ncondes/go-social/commit/9bf02fc33d1627e6939589b961bfd1f70ea57f59))
* integrate swagger documentation generation and update api post endpoint for versioning ([7cb0d23](https://github.com/ncondes/go-social/commit/7cb0d23f3e5d58cbbb4454c3742b7f23f9a8bd7a))
* post crud operations and post and get comments ([33bb42b](https://github.com/ncondes/go-social/commit/33bb42b4fd13cf1daa645ec8cf2d47e0f52910e0))
* scaffolding api server ([11b1f88](https://github.com/ncondes/go-social/commit/11b1f88aec224ac15ebf5dbd923609776d7f1efb))
* user feed implementation with score algorithm and filtering ([edc7783](https://github.com/ncondes/go-social/commit/edc7783339728556066917034bf75b1f79771917))


### Bug Fixes

* add swagger documentation generation to ci pipeline ([95032d4](https://github.com/ncondes/go-social/commit/95032d4949c0ba11d4a68830b845ba2ebace087a))
* remove duplicate db package import ([bef2bae](https://github.com/ncondes/go-social/commit/bef2baedc1b853efe6643e03d06bacec1d827b05))
* reorder swagger generation before build ([4639846](https://github.com/ncondes/go-social/commit/46398469bce641a35a3c21e042476c78b2ecfc84))
