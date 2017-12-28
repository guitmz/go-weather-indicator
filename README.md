# go-weather-indicator
Minimalistic GTK weather indicator written in GoLang. Uses [Weatherbit](https://www.weatherbit.io) as provider. 

You can get yout free API key here: https://www.weatherbit.io/account/create

# Installation
Either download the binary release from GitHub or run:
`$ go get -u "github.com/guitmz/go-weather-indicator"`

# Usage
`$ go-weather-indicator --city Berlin --country Germany --key API_KEY`

# TODO
- More weather information
- Better error handling
- Try to cleanup and add comments in the code
- Write tests
- Allow more customization like displaying temeprature in Farenheit
