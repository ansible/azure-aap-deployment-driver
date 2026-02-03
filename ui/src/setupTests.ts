// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom';

// Polyfills for Web APIs used by react-router v7 in Jest environment
import { TextEncoder, TextDecoder } from 'util';
// @ts-ignore
if (!global.TextEncoder) {
	// @ts-ignore
	global.TextEncoder = TextEncoder;
}
// @ts-ignore
if (!global.TextDecoder) {
	// @ts-ignore
	global.TextDecoder = TextDecoder as unknown as typeof global.TextDecoder;
}