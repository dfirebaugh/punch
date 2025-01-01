export const renderJSON = (data, container) => {
  Object.entries(data).forEach(([key, value]) => {
    const containerItem = document.createElement('div');
    containerItem.classList.add('json-item');
    containerItem.classList.add('collapsed');

    const toggleButton = document.createElement('button');
    toggleButton.classList.add('toggle-button');
    toggleButton.textContent = '+';
    toggleButton.addEventListener('click', () => {
      const isCollapsed = containerItem.classList.contains('collapsed');
      toggleAllChildren(containerItem, isCollapsed);
      toggleButton.textContent = isCollapsed ? '-' : '+';
    });

    const keyElement = document.createElement('span');
    keyElement.classList.add('json-key');
    keyElement.textContent = `"${key}": `;
    keyElement.addEventListener('click', () => {
      const isCollapsed = containerItem.classList.toggle('collapsed');
      toggleButton.textContent = isCollapsed ? '+' : '-';
    });


    const valueElement = document.createElement('div');
    valueElement.classList.add('json-value');

    const summaryElement = document.createElement('span');
    summaryElement.classList.add('json-summary');
    summaryElement.style.color = '#aaa';
    summaryElement.style.fontSize = '0.9em';
    summaryElement.style.marginLeft = '8px';
    summaryElement.style.display = 'inline';

    if (typeof value === 'object' && value !== null) {
      const childCount = Object.keys(value).length;
      summaryElement.textContent = `(${childCount} ${childCount === 1 ? 'node' : 'nodes'})`;

      const bracketOpen = document.createElement('span');
      bracketOpen.classList.add('json-bracket');
      bracketOpen.textContent = '{';

      const bracketClose = document.createElement('span');
      bracketClose.classList.add('json-bracket');
      bracketClose.textContent = '}';

      valueElement.appendChild(bracketOpen);
      const nestedContainer = document.createElement('div');
      nestedContainer.style.marginLeft = '20px';
      renderJSON(value, nestedContainer);
      valueElement.appendChild(nestedContainer);
      valueElement.appendChild(bracketClose);

      summaryElement.style.display = 'inline';
    } else {
      valueElement.textContent = JSON.stringify(value, null, 2);
      summaryElement.textContent = '';
    }

    containerItem.appendChild(toggleButton);
    containerItem.appendChild(keyElement);
    containerItem.appendChild(summaryElement);
    containerItem.appendChild(valueElement);

    container.appendChild(containerItem);
  });
};

const toggleAllChildren = (element, expand) => {
  const toggleNodes = element.querySelectorAll('.json-item');
  toggleNodes.forEach((child) => {
    if (expand) {
      child.classList.remove('collapsed');
      element.classList.remove('collapsed');
    } else {
      child.classList.add('collapsed');
      element.classList.add('collapsed');
    }
    const button = child.querySelector('.toggle-button');
    if (button) button.textContent = expand ? '-' : '+';
  });
};

