/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./internal/**/*.templ'],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

