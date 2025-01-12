import resolve from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import typescript from "@rollup/plugin-typescript";
import postcss from 'rollup-plugin-postcss';
import peerDepsExternal from 'rollup-plugin-peer-deps-external';
import { babel } from '@rollup/plugin-babel';
import packageJson from "./package.json" assert { type: 'json' };

console.log()
export default [
    {
        input: "src/index.ts",
        output: [
            {
                file: packageJson.main,
                format: "esm",
                sourcemap: true,
            },
        ],
        plugins: [
            peerDepsExternal(),
            resolve({
                browser: true,
                preferBuiltins: true
            }),
            commonjs(),
            typescript({ tsconfig: "./tsconfig.json" }),
            postcss(),
            babel({ babelHelpers: 'runtime', exclude: 'node_modules/**' }),
        ],
        external: ["react", "react-dom", "react/jsx-runtime"],
    },
];