import vue from 'rollup-plugin-vue';
import buble from 'rollup-plugin-buble';
export default {
    input: './build/module.js',
    output: {
        name: 'HelloWorld',
        exports: 'named',
    },
    plugins: [
        vue({
            css: true,
            compileTemplate: true,
        }),
        buble(),
    ],
};