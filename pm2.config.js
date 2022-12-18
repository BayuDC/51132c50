module.exports = {
    apps: [
        {
            name: 'tink',
            script: '.out/tink',
            env: {
                PORT: 8080,
                GIN_MODE: 'release',
            },
        },
    ],
};
