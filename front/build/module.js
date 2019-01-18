// Import vue components
import HelloWorld from '../src/components/HelloWorld.vue';

const components = {
	HelloWorld
};

// Declare install function executed by Vue.use()
export default function install(Vue) {
	if (install.installed) return;
	install.installed = true;
	Object.keys(components).forEach(name => {
		Vue.component(name, components[name]);
	});
}
