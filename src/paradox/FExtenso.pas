unit FExtenso;

interface

Function Extenso(Valor: double):string;

implementation

uses SysUtils;

const
Centenas: array[1..9] of string[12]=('CEM','DUZENTOS','TREZENTOS','QUATROCENTOS',
                                      'QUINHENTOS','SEISCENTOS','SETECENTOS',
                                      'OITOCENTOS','NOVECENTOS');
Dezenas : array[2..9] of string[10]=('VINTE','TRINTA','QUARENTA','CINQUENTA',
                                      'SESSENTA','SETENTA','OITENTA','NOVENTA');
Dez     : array[0..9] of string[10]=('DEZ','ONZE','DOZE','TREZE','QUATORZE',
                                      'QUINZE','DEZESSEIS','DEZESSETE',
                                      'DEZOITO','DEZENOVE');
Unidades: array[1..9] of string[10]=('UM','DOIS','TRES','QUATRO','CINCO',
                                      'SEIS','SETE','OITO','NOVE');
 MoedaSingular = 'REAL';
 MoedaPlural   = 'REAIS';
 CentSingular  = 'CENTAVO';
 CentPlural    = 'CENTAVOS';
 Zero          = 'ZERO';

Function Ext3(Parte:string):string;
var
 Base: string;
 digito: integer;
begin
Base:='';
digito:=StrToInt(Parte[1]);
if digito=0 then
 Base:=''
else
 Base:=Centenas[digito];
if (digito = 1) and (Parte > '100') then
 Base:='CENTO';
Digito:=StrToInt(Parte[2]);
if digito = 1 then
 begin
  Digito:=StrToInt(Parte[3]);
  if Base <> '' then
   Base:=Base + ' E ';
  Base:=Base + Dez[Digito];
 end
else
 begin
  if (Base <> '') and (Digito > 0) then
   Base:=Base+' E ';
  if Digito > 1 then
   Base:=Base + Dezenas[digito];
  Digito:=StrToInt(Parte[3]);
  if Digito > 0 then
   begin
    if Base <> '' then Base:=Base+' E ';
    Base:=Base+Unidades[Digito];
   end;
 end;
Result:=Base;
end;

Function Extenso;
var
 ComoTexto: string;
 Parte: string;
begin
Result:='';
ComoTexto:=FloatToStrF(Abs(Valor),ffFixed,18,2);
// Acrescenta zeros a esquerda ate 12 digitos
while length(ComoTexto) < 15 do
 Insert('0',ComoTexto,1);
// Retira caracteres a esquerda para o máximo de 12 digitos
while length(ComoTexto) > 15 do
 delete(ComoTexto,1,1);

// Calcula os bilhões
Parte:=Ext3(copy(ComoTexto,1,3));
if StrToInt(copy(ComoTexto,1,3)) = 1 then
 Parte:=Parte + ' BILHAO'
else
 if Parte <> '' then
  Parte:=Parte + ' BILHOES';
Result:=Parte;

// Calcula os nilhões
Parte:=Ext3(copy(ComoTexto,4,3));
if Parte <> '' then
 begin
  if Result <> '' then
   Result:=Result+', ';
  if StrToInt(copy(ComoTexto,4,3)) = 1 then
   Parte:=Parte + ' MILHAO'
  else
   Parte:=Parte + ' MILHOES';
  Result:=Result + Parte;
 end;





interface

uses
  Windows, Messages, SysUtils, Variants, Classes, Graphics, Controls, Forms,
  Dialogs, QuickRpt, DB, DBTables, QRCtrls, ExtCtrls;

type
  TForm9 = class(TForm)
  private
    { Private declarations }
  public
    { Public declarations }
  end;

var
  Form9: TForm9;

implementation
   uses udmorca;
{$R *.dfm}

end.
