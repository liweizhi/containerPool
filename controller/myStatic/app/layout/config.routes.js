(function() {
    'use strict';

    angular
        .module('container-pool.layout')
        .config(getRoutes);

    getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];

    function getRoutes($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('dashboard', {
                url: '/',
                abstract: true,
                templateUrl: 'app/layout/dashboard.html',
            });
    }
})();

