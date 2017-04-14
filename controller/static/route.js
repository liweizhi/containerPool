/**
 * Created by james on 17/2/23.
 */

var myApp = angular.module('containerPool', ['ui.router']);

myApp.config(function($stateProvider) {
    var containersState = {
        name: 'containers',
        url: '',
        template: '<h3 class="starter-template">containers</h3>'
    }

    var imagesState = {
        name: 'images',
        url: '/images',
        template: '<h3 class="starter-template">images</h3>'
    }
    var accountsState = {
        name: 'accounts',
        url: '/accounts',
        template: '<h3 class="starter-template">accounts</h3>'
    }
    var nodesState = {
        name: 'nodes',
        url: '/nodes',
        template: '<h3 class="starter-template">nodes</h3>'
    }

    var helloState = {
        name: 'hello',
        url: '/hello',
        template: '<h3>hello world!</h3>'
    }

    var aboutState = {
        name: 'about',
        url: '/about',
        template: '<h3>Its the UI-Router hello world app!</h3>'
    }

    $stateProvider.state(helloState);
    $stateProvider.state(aboutState);
    $stateProvider.state(containersState);
    $stateProvider.state(imagesState);
    $stateProvider.state(accountsState);
    $stateProvider.state(nodesState);


});