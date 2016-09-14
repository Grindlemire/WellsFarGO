var NODE_ENV = process.env.NODE_ENV || JSON.stringify(process.env.NODE_ENV || 'development');
var watch = process.env.WATCH == "true";
var debug = NODE_ENV == 'development';

var autoprefixer = require('autoprefixer'); // Don't worry about vendor prefixing like webkit.
var ExtractTextPlugin = require('extract-text-webpack-plugin'); // Extracts css into its own bundle (faster loading of css)
var CopyWebpackPlugin = require('copy-webpack-plugin'); // Copies files and dirs from dev to a build dir

var path = require('path');
var webpack = require('webpack');

module.exports = {
    watch: watch ? true : false,
    context: __dirname,
    devtool: debug ? "cheap-module-source-map" : "",
    entry: {
        index: __dirname + "/rest/dev/js/index.js",
    },
    output: {
        path: __dirname + "/rest/build/",
        filename: "js/[name].min.js"
    },
    module: {
        noParse: [
        ],
        loaders: [
            {
                test: /\.scss$/,
                exclude: /(node_modules)/,
                loader: ExtractTextPlugin.extract('style-loader', 'css-loader!postcss-loader!sass-loader')
            },
            {
                test: /\.css$/,
                loader: 'style-loader!css-loader'
            },
            {
               test: /\.(png|jpg|jpeg|gif|svg|woff|woff2|ttf|eot)$/,
               loader: 'file?name=stylesheets/libs/[name].[ext]',
            },
            {
                test: /\.js?$/,
                exclude: /(node_modules)/,
                loader: 'babel-loader',
                query: {
                  presets: ['react', 'es2015', 'stage-2'],
                  plugins: ['transform-class-properties']
                }
            }
        ]
    },
    postcss: debug ? [] : [
        autoprefixer({
           browsers: ['last 2 version']
        }),
        require('cssnano')
    ],
    plugins: debug ? [
        new ExtractTextPlugin("styles/[name].min.css"),
        new CopyWebpackPlugin([
            {from: './rest/dev/static/', to: "static"},
        ])
    ] : [
        new ExtractTextPlugin("styles/[name].min.css"),
        new CopyWebpackPlugin([
            {from: './rest/dev/static/', to: "static"},
        ]),


        new webpack.DefinePlugin(
            {
                'process.env': {
                    'NODE_ENV': JSON.stringify('production')
                }
            }),
        new webpack.optimize.DedupePlugin(),
        new webpack.optimize.OccurenceOrderPlugin(),
        new webpack.optimize.UglifyJsPlugin({ mangle: false, sourcemap: false, compress: {warnings: false}}),
    ],
    resolve: {
        extensions: ['', '.js', '.scss', '.css'],
        root: [path.resolve('./rest/dev')]
    }
};
