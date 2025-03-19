module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  parserOptions: {
    project: 'tsconfig.json',
    tsconfigRootDir: __dirname,
    sourceType: 'module'
  },
  plugins: ['@typescript-eslint', 'react-refresh', '@tanstack/query'],
  extends: [
    'eslint:recommended',
    'plugin:react/recommended',
    'plugin:react-hooks/recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:@typescript-eslint/recommended-requiring-type-checking'
  ],
  overrides: [
    {
      files: ['*.js', '*.jsx', '*.ts', '*.tsx']
    }
  ],
  ignorePatterns: [
    '.eslintrc.cjs',
    'jest.config.js',
    'build.js',
    'babel.config.js',
    'tests/**',
    'node_modules/**',
    'dist/**',
    'build/**',
    'tests-reports/**',
    'src/config/**',
    'postcss.config.js',
    'vite.config.ts',
    'tailwind.config.js',
    'utils/**',
  ],
  settings: {
    react: {
      version: 'detect'
    }
  },
  rules: {
    'block-spacing': ['warn', 'always'],
    'brace-style': ['warn', '1tbs'],
    curly: ['warn', 'all'],
    indent: [
      'warn',
      2,
      {
        SwitchCase: 1,
        ignoredNodes: ['JSXElement *', 'JSXElement']
      }
    ],
    'no-mixed-spaces-and-tabs': 'error',
    'no-multi-spaces': 'warn',
    'no-trailing-spaces': 'warn',
    'no-whitespace-before-property': 'error',
    'react/jsx-boolean-value': 'off',
    'react/jsx-indent-props': ['off'],
    'space-before-blocks': ['warn', 'always'],
    'space-before-function-paren': [
      'warn',
      {
        named: 'never',
        anonymous: 'ignore',
        asyncArrow: 'always'
      }
    ],
    'space-in-parens': ['warn', 'never'],
    '@typescript-eslint/no-explicit-any': 'off',
    '@typescript-eslint/no-inferrable-types': 'off',
    '@typescript-eslint/no-unsafe-assignment': 'off',
    '@typescript-eslint/no-unsafe-member-access': 'off',
    '@typescript-eslint/no-unsafe-return': 'off',
    '@typescript-eslint/no-misused-promises': 'off',
    '@typescript-eslint/no-unused-vars': ['warn', {argsIgnorePattern: '^_'}],
    '@typescript-eslint/restrict-template-expressions': 'off',
    '@typescript-eslint/unbound-method': 'off',
    '@typescript-eslint/no-unsafe-call': 'warn',
    'no-async-promise-executor': 'warn',
    'react/prop-types': 'off',
    'react/no-unknown-property': 'warn',
    '@typescript-eslint/no-redundant-type-constituents': 'off',
    "@tanstack/query/exhaustive-deps": "error",
    "@tanstack/query/no-rest-destructuring": "warn",
    "@tanstack/query/stable-query-client": "error",
    'react/react-in-jsx-scope': 'off',
    'max-len': 'off'
  }
};
