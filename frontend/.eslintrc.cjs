module.exports = {
  root: true,
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: ["eslint:recommended", "plugin:vue/vue3-recommended"],
  parser: "vue-eslint-parser",
  parserOptions: {
    ecmaVersion: "latest",
    sourceType: "module",
    parser: "espree",
  },
  rules: {
    "no-unused-vars": "warn",
    "no-console": "warn",
    "vue/multi-word-component-names": "off",
    "vue/no-unused-vars": "warn",
  },
  ignorePatterns: ["dist/", "node_modules/"],
};
