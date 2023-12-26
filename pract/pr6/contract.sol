// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract security_paper_exchange {
    // Структура для представления ценной бумаги
    struct security_paper {
        uint id; // Уникальный идентификатор ценной бумаги
        string name; // Название ценной бумаги
        uint price; // Цена ценной бумаги
        address payable owner; // Адрес владельца ценной бумаги
        address payable buyer; // Адрес покупателя ценной бумаги
    }

    uint public total_securities = 0; // Общее количество ценных бумаг
    mapping(uint => security_paper) private securitiesList; // Хранение ценных бумаг по их идентификатору

    // Модификаторы для проверки прав доступа и условий
    modifier onlyOwner(uint _security_paper_id) {
        require(
            msg.sender == securitiesList[_security_paper_id].owner,
            "Only the owner can perform this action"
        );
        _;
    }

    modifier notAlreadySold(uint _security_paper_id) {
        require(
            securitiesList[_security_paper_id].buyer == payable(address(0)),
            "This security paper has already been sold"
        );
        _;
    }

    modifier enoughPrice(uint _security_paper_id) {
        require(
            msg.value >= uint(securitiesList[_security_paper_id].price),
            "Insufficient funds to buy this security paper"
        );
        _;
    }

    // Добавление новой ценной бумаги
    function add_security_paper(
        string memory _security_paper_name,
        uint _price
    ) public {
        require(msg.sender != address(0), "Invalid address");
        total_securities++;

        // Создание новой ценной бумаги и добавление ее в список
        securitiesList[total_securities] = security_paper(
            total_securities, // Уникальный идентификатор увеличивается на каждую новую бумагу
            _security_paper_name,
            _price,
            payable(msg.sender), // Владелец по умолчанию — адрес отправителя транзакции
            payable(address(0)) // Начальный покупатель отсутствует
        );
    }

    // Покупка ценной бумаги
    function buy_security_paper(uint _security_paper_id)
        public
        payable
        enoughPrice(_security_paper_id)
        notAlreadySold(_security_paper_id)
    {
        require(msg.sender != address(0), "Invalid address");
        address payable _owner = securitiesList[_security_paper_id].owner;
        uint totalCost = securitiesList[_security_paper_id].price;

        // Перевод оплаты владельцу ценной бумаги
        _owner.transfer(totalCost);

        // Установка покупателя ценной бумаги
        securitiesList[_security_paper_id].buyer = payable(msg.sender);
    }

    // Обмен владельца ценной бумаги
    function exchange_security_paper_ownership(
        uint _security_paper_id,
        address payable _newOwner
    )
        public
        onlyOwner(_security_paper_id)
    {
        require(_newOwner != address(0), "Invalid address");

        // Обновление владельца ценной бумаги
        securitiesList[_security_paper_id].owner = _newOwner;
        // Сброс адреса покупателя при изменении владельца
        securitiesList[_security_paper_id].buyer = payable(address(0));
    }

    // Получение информации о ценной бумаге
    function get_security_paper_details(uint _security_paper_id)
    public
    view
    returns (
        string memory,
        uint,
        address,
        address
    )
    {
        // Проверка, что ценная бумага была продана или запрос от владельца
        require(
            securitiesList[_security_paper_id].buyer != payable(address(0)) || msg.sender == securitiesList[_security_paper_id].owner,
            "Only the owner can perform this action"
        );

        // Возвращение данных о ценной бумаге: название, цена, владелец, покупатель
        return (
            securitiesList[_security_paper_id].name,
            securitiesList[_security_paper_id].price,
            securitiesList[_security_paper_id].owner,
            securitiesList[_security_paper_id].buyer
        );
    }

}
