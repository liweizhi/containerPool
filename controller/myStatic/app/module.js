(function(){
    'use strict';

    angular
        .module('container-pool', [
                'container-pool.toplogy',
                'container-pool.accounts',
                'container-pool.core',
                'container-pool.services',
                'container-pool.layout',
                'container-pool.login',
                'container-pool.containers',
                'container-pool.images',
                'container-pool.nodes',
                'container-pool.filters',
                'angular-jwt',
                'base64',
                'selectize',
                'ui.router'
        ]);

})();
