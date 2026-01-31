export const apiFetch = async (url, options = {}) => {
  const headers = {
    ...options.headers,
    "ngrok-skip-browser-warning": "69420",
  };

  return fetch(url, {
    ...options,
    headers,
  });
};

export const enhanceTask = async (taskId, token) => {
    return apiFetch(`/api/tasks/enhance/${taskId}`, {
        method: "POST",
        headers: {
            Authorization: `Bearer ${token}`
        }
    });
};
