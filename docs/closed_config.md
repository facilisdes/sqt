# Закрытый конфиг (.sqt)
Файл используется только сервером. Создаётся путём заполнения базового конфига и последующего шифрования (запаковки).  
Пример базового файла:
```bash
read_time_init=300
read_time_step=1000
max_stack_size=10
read_time_growth=sum
read_time_parameter=0
```
## Описание параметров
- **read_time_init** - начальное время чтения в очереди. Первый запрос в очереди будет исполняться не быстрее этого времени.  
- **read_time_step** - шаг времени чтения в очереди. Отвечает за приращение времени чтения для следующих за первым элементов очереди.  
- **max_stack_size** - максимальный размер очереди. При достижении размера очереди, равному значению этого параметра, запросы отклоняются.  
- **read_time_growth** - функция роста времени чтения. Возможные значения и описания функций приводятся ниже.  
- **read_time_parameter** - параметр функции роста времени чтения.   
      
##Функции роста времени чтения
Во всех формулах n - количество элементов в очереди на момент поступления запроса.

<table>
      <thead>
            <tr>
                  <td>Функция</td>
                  <td>Формула</td>
                  <td>Описание</td>
            </tr>
      </thead>
      <tbody>
            <tr>
                  <td>**sum**</td>
                  <td>*read_time_init + read_time_step * n*</td>
                  <td>Время растёт линейно. </td>
            </tr>
            <tr>
                  <td>**msum**</td>
                  <td>*read_time_init + (read_time_step * n * read_time_parameter)*</td>
                  <td>Время растёт линейно. Зависимость от размера очереди регулируется через параметр *read_time_parameter*</td>
            </tr>
            <tr>
                  <td>**exp**</td>
                  <td>*read_time_init + read_time_step ^ n)*</td>
                  <td>Быстрый нелинейный рост.</td>
            </tr>
            <tr>
                  <td>**mexp**</td>
                  <td>*read_time_init + read_time_step ^ (n * read_time_parameter)*</td>
                  <td>Быстрый нелинейный рост. Зависимость от размера очереди регулируется через параметр *read_time_parameter*</td>
            </tr>
            <tr>
                  <td>*log*</td>
                  <td>*read_time_init + read_time_step * log(1 + n)*</td>
                  <td>Медленный нелинейный рост.</td>
            </tr>
            <tr>
                  <td>*mlog*</td>
                  <td>*read_time_init + read_time_step * log((1 + n) * read_time_parameter)* </td>
                  <td>Медленный нелинейный рост. Зависимость от размера очереди регулируется через параметр *read_time_parameter*</td>
            </tr>
      </tbody>
</table>
