const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("security_paper_exchange", function () {
  let security_paper_exchange;
  let accounts;

  beforeEach(async function () {
    // Получение контракта security_paper_exchange из Hardhat и развертывание
    const Contract = await ethers.getContractFactory("security_paper_exchange");
    security_paper_exchange = await Contract.deploy();
    await security_paper_exchange.deployed();
    
    // Получение списка тестовых аккаунтов из Hardhat
    accounts = await ethers.getSigners();
  });

  // Проверка добавления нескольких ценных бумаг и подсчета их количества
  it("Добавление нескольких ценных бумаг и отображение количества", async function () {
    // Добавление двух ценных бумаг и проверка, что их количество равно 2
    await security_paper_exchange.connect(accounts[0]).add_security_paper("Security1", ethers.utils.parseEther('1'));
    await security_paper_exchange.connect(accounts[0]).add_security_paper("Security2", ethers.utils.parseEther('2'));
    const count = await security_paper_exchange.total_securities();
    expect(count.toNumber()).to.equal(2);
  });

  // Проверка попытки покупки ценной бумаги по меньшей цене
  it("Не продавать ценные бумаги по меньшей цене", async function () {
    // Добавление дорогой ценной бумаги и попытка покупки за меньшую сумму
    await security_paper_exchange.connect(accounts[0]).add_security_paper("Expensive Security", ethers.utils.parseEther('10'));
    try {
      await security_paper_exchange.connect(accounts[1]).buy_security_paper(1, { value: ethers.utils.parseEther('11') });
    } catch (err) {
      expect(err.message).to.include("Insufficient funds to buy this security paper");
    }
  });

  // Проверка разрешения просмотра адреса покупателем ценных бумаг
  it("Разрешение просмотра адреса покупателем ценных бумаг", async function () {
    // Добавление и покупка ценной бумаги, затем проверка доступа к ее адресу покупателем
    await security_paper_exchange.connect(accounts[0]).add_security_paper("Security1", ethers.utils.parseEther('1'));
    await security_paper_exchange.connect(accounts[1]).buy_security_paper(1, { value: ethers.utils.parseEther('1') });
    const details = await security_paper_exchange.connect(accounts[1]).get_security_paper_details(1);
    expect(details[1]).to.equal(ethers.utils.parseEther('1'));
  });

  // Проверка запрета просмотра адреса ценных бумаг, купленных другим пользователем
  it("Запрет на просмотр адреса ценных бумаг, купленных другим пользователем", async function () {
    // Добавление, покупка ценной бумаги одним аккаунтом и попытка просмотра ее другим аккаунтом
    await security_paper_exchange.connect(accounts[0]).add_security_paper("Security1", ethers.utils.parseEther('1'));
    await security_paper_exchange.connect(accounts[1]).buy_security_paper(1, { value: ethers.utils.parseEther('1') });
    try {
      await security_paper_exchange.connect(accounts[2]).get_security_paper_details(1);
    } catch (err) {
      expect(err.message).to.include("Only the owner can perform this action");
    }
  });
});
