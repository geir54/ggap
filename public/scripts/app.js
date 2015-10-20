app = angular.module('app', ['ngRoute']);

app.config(function($routeProvider, $locationProvider){
  $routeProvider
    .when('/', {templateUrl: '/partials/main.html'})
    .when('/login', {templateUrl: '/partials/login.html',
                   controller: 'AuthenticationController'})
    .when('/signup', {templateUrl: '/partials/signup.html',
                   controller: 'AuthenticationController'})
    .otherwise({redirectTo: '/'});

  $locationProvider.html5Mode(true);
});
