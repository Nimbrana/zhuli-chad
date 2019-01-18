const gulp = require('gulp')
const replace = require('gulp-string-replace')
const del = require('del');

const filename = 'dist/HelloWorld.js'

gulp.task('clean', function () {
  return del([
    'dist/**/*',
  ]);
});

gulp.task('copy', () => gulp.src(['./build/package.json']).pipe(gulp.dest('./dist')))

gulp.task('optimize', function() {
  return gulp
    .src([filename])
    .pipe(replace(', createInjectorSSR', ''))
    .pipe(gulp.dest('dist'))
})

gulp.task('package', gulp.series('optimize', 'copy'))
