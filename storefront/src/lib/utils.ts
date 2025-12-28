export const getPlaceholderImage = (text: string = 'No Image') => {
    // Simple SVG placeholder data URI
    const svg = `
    <svg width="400" height="400" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid slice" focusable="false" role="img" aria-label="Placeholder">
        <rect width="100%" height="100%" fill="#f3f4f6"></rect>
        <text x="50%" y="50%" fill="#9ca3af" dy=".3em" font-family="sans-serif" font-size="24" text-anchor="middle">${text}</text>
    </svg>`;
    return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`;
};
