(function(){
	'use strict';

	angular
		.module('container-pool.images')
		.config(getRoutes);

	getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];

	function getRoutes($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('dashboard.images', {
			    url: '^/images',
			    templateUrl: 'app/images/images.html',
                            controller: 'ImagesController',
                            controllerAs: 'vm',
                            authenticate: 'true',
                            resolve: {
                                images: ['ImagesService', '$state', '$stateParams', function (ImagesService, $state, $stateParams) {
                                    return ImagesService.list().then(null, function(errorData) {	                            
                                        $state.go('error');
                                    }); 
                                }] 
                            }
			});
	}
})();
