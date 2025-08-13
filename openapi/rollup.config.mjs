import copy from 'rollup-plugin-copy';

export default {
  input: './index.js',
  plugins: [
    copy({
      targets: [
        { src: 'node_modules/redoc/bundles/redoc.standalone.js', dest: 'public/docs/' },
        {
          src: [
            'node_modules/swagger-ui-dist/swagger-ui-bundle.js',
            'node_modules/swagger-ui-dist/swagger-ui-standalone-preset.js',
            'node_modules/swagger-ui-dist/swagger-ui.css',
            'node_modules/swagger-ui-dist/favicon-32x32.png',
            'node_modules/swagger-ui-dist/favicon-16x16.png',
          ], dest: 'public/swagger/'
        },
      ],
    }),
  ],
};
