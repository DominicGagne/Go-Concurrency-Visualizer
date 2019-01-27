var scene, camera, renderer, controls, orbit, trace, cameraControls;

init();
Leap.loop(function(frame){
  cameraControls.update(frame); // rotating, zooming & panning
  trace.onControlsChanged(camera);
});

animate();

function init() {
	var scale = 1;

	// STATS
	stats = new Stats();
	stats.setMode( 0 ); // 0: fps, 1: ms, 2: mb
	stats.domElement.style.position = 'absolute';
	stats.domElement.style.left = '0px';
	stats.domElement.style.top = '0px';
	document.body.appendChild( stats.domElement );

	controller = new Leap.Controller();
	scene = new THREE.Scene();

	// CAMERA
	width = window.innerWidth;
	height = window.innerHeight;
	var center = new THREE.Vector3(60, -50, -10);
	//camera = new THREE.OrthographicCamera( width / - 2, width / 2, height / 2, height / - 2, -1000*scale, 2000*scale );
	camera = new THREE.PerspectiveCamera(75, width / height, 1, 1000 * scale );
	camera.position.z = 100 * scale;
	camera.position.y = 50 * scale;
	camera.updateProjectionMatrix();
	
	mat1 = new THREE.LineBasicMaterial( { color: 0x0000ff, linewidth: 4, } );
	trace = new GoThree.Trace();
	trace.init(scene, data, params, scale);

	// RENDERER
	renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true, });
	renderer.setSize( width, height );
	renderer.setClearColor( '#1D1F17', 1);

	// light for hand
	var light = new THREE.AmbientLight( 0x505050 ); // soft white light
	scene.add( light );

	// leap camera controls
	//controls = new THREE.LeapMyControls( camera , controller, renderer.domElement );
	//controls = new THREE.LeapPointerControls( camera , controller, renderer.domElement );
	cameraControls = new THREE.LeapCameraControls(camera);
	cameraControls.panEnabled      = true;
	cameraControls.panSpeed        = 1.0;
	cameraControls.panHands        = 2;
	cameraControls.panFingers      = [6,12];
	cameraControls.panRightHanded  = true; // right-handed
	cameraControls.panHandPosition = true; // palm position used
	cameraControls.panStabilized   = true; // stabilized palm position used
	
	cameraControls.rotateEnabled         = true;
	cameraControls.rotateHands           = 1;
	cameraControls.rotateSpeed           = 0.8;
	cameraControls.rotateFingers         = [4, 5];
	cameraControls.rotateRightHanded     = true;
	cameraControls.rotateHandPosition    = false;
	cameraControls.rotateStabilized      = true;

	cameraControls.zoomEnabled         = false;
	cameraControls.zoomHands           = 2;
	cameraControls.zoomSpeed           = 1;
	cameraControls.zoomFingers         = [6, 12];
	cameraControls.zoomRightHanded     = true;
	cameraControls.zoomHandPosition    = true;
	cameraControls.zoomStabilized      = true;

	// CONTROLS
	orbit = new THREE.OrbitControls( camera, renderer.domElement );
	orbit.autoRotate = false;
	orbit.autoRotateSpeed = 1.0;
	orbit.addEventListener( 'change', function() {
		trace.onControlsChanged(orbit.object);
	});

	// ADD CUSTOM KEY HANDLERS
	document.addEventListener("keydown", function(event) {keydown(event)}, false);

	console.log("slowing down trace")
	trace.slowdown();
	trace.slowdown();
	trace.slowdown();
	trace.slowdown();
	trace.slowdown();


	document.body.appendChild( renderer.domElement );

	controller.connect();
}

function animate() {
	if (orbit.autoRotate) {
		orbit.update();
	};
	trace.animate();

	requestAnimationFrame(animate);
	stats.begin();
	renderer.render(scene, camera);
    stats.end();
}

function keydown(event) {
	console.log(event.which);
	switch (event.which) {
		case 80: // 'P' - (Un)Pause autoRotate
			toggleAutoRotate();
			break;
		case 82: // 'R' - Reset
			trace.resetTime();
			break;
		case 83: // 'S' - Slower
			trace.slowdown();
			break;
		case 70: // 'F' - Faster
			trace.speedup();
			break;
		case 187: // '+' - IncWidth
			trace.incWidth();
			break;
		case 189: // '-' - DecWidth
			trace.decWidth();
			break;
		case 48: // '0' - HighlighMode Default
			trace.highlight("default");
			break;
		case 49: // '1' - HighlighMode 1
			trace.highlight("mode1");
			break;
		case 50: // '2' - HighlighMode 2
			trace.highlight("mode2");
			break;
		case 51: // '2' - HighlighMode 3
			trace.highlight("mode3");
			break;
		case 52: // '2' - HighlighMode 4
			trace.highlight("mode4");
			break;
	}
}

function toggleAutoRotate() {
	orbit.autoRotate = !orbit.autoRotate;
}
